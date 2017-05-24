package linda

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
	neturl "net/url"
	"time"
)

type RedisBroker struct {
	redisURL *neturl.URL
	pool     *redis.Pool
}

func (r *RedisBroker) Connect(url *neturl.URL) error {
	var host = url.Host
	var password string
	var db string
	if url.User != nil {
		password, _ = url.User.Password()
	}
	if len(url.Path) > 1 {
		db = url.Path[1:]
	}
	r.pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host)
			if err != nil {
				return nil, err
			}
			if password != "" {
				_, err := c.Do("AUTH", password)
				if err != nil {
					c.Close()
					return nil, err
				}
			}
			if db != "" {
				_, err := c.Do("SELECT", db)
				if err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, nil
		},
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return nil
}

func (r *RedisBroker) Close() error {
	return r.pool.Close()
}

func (r *RedisBroker) MigrateExpiredJobs(queue string) {
	r.migrateExpiredJobs(fmt.Sprintf("%s:reserved", queue), queue)
	r.migrateExpiredJobs(fmt.Sprintf("%s:delayed", queue), queue)
}

func (r *RedisBroker) migrateExpiredJobs(from string, to string) {
	conn := r.pool.Get()
	defer conn.Close()
	logrus.Debugf("migrate expired jobs from %s to %s", from, to)
	_, err := conn.Do("EVAL", MigrateJobsScript, 2, from, to, time.Now().Unix())
	if err != nil {
		logrus.Error(err)
	}
}
func (r *RedisBroker) Pop(queue string, ack bool, timeout int64) (job *Job, err error) {
	conn := r.pool.Get()
	defer conn.Close()
	// reserve next job
	var reply []byte
	if ack {
		reply, err = redis.Bytes(conn.Do("EVAL", PopJobScript, 2, queue, fmt.Sprintf("%s:reserved", queue), delayAt(timeout)))
	} else {
		reply, err = redis.Bytes(conn.Do("LPOP", queue))
	}
	if err == redis.ErrNil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(reply, &job)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	logrus.Debugf("pop job {%s}", job)
	return job, nil
}

func (r *RedisBroker) DeleteReserved(queue string, job *Job) error {
	conn := r.pool.Get()
	defer conn.Close()
	bytes, err := json.Marshal(job)
	if err != nil {
		logrus.Error(err)
		return err
	}
	_, err = conn.Do("ZREM", fmt.Sprintf("%s:reserved", queue), bytes)
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Infof("delete reserved job {%s}", job)
	return nil
}

func (r *RedisBroker) DeleteAndRelease(queue string, job *Job, delay int64) error {
	conn := r.pool.Get()
	defer conn.Close()
	bytes, err := json.Marshal(job)
	if err != nil {
		logrus.Error(err)
		return err
	}
	_, err = conn.Do("EVAL", ReleaseScript, 2, fmt.Sprintf("%s:delayed", queue), fmt.Sprintf("%s:reserved", queue), bytes, delayAt(delay))
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Infof("delete and release job {%s}", job)
	return nil
}

func (r *RedisBroker) Push(job *Job, queue string) error {
	conn := r.pool.Get()
	defer conn.Close()
	bytes, err := json.Marshal(job)
	if err != nil {
		logrus.Error(err)
		return err
	}
	_, err = conn.Do("RPUSH", queue, bytes)
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Infof("push job {%s}", job)
	return nil
}

func (r *RedisBroker) Later(delay int64, job *Job, queue string) error {
	conn := r.pool.Get()
	defer conn.Close()
	bytes, err := json.Marshal(job)
	if err != nil {
		logrus.Error(err)
		return err
	}
	_, err = conn.Do("ZADD", fmt.Sprintf("%s:delayed", queue), delayAt(delay), bytes)
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Infof("push job {%s} with delay %d", job, delay)
	return nil
}

func delayAt(delay int64) int64 {
	return time.Now().Unix() + delay
}

package linda

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
	neturl "net/url"
	"time"
)

// RedisBroker  broker driver with redis
type RedisBroker struct {
	redisURL *neturl.URL
	pool     *redis.Pool
}

// Connect broker backend with url
func (r *RedisBroker) Connect(url *neturl.URL) error {
	r.redisURL = url

	var network string
	var host string
	var password string
	var db string
	network = "tcp"
	host = url.Host
	if url.User != nil {
		password, _ = url.User.Password()
	}
	if len(url.Path) > 1 {
		db = url.Path[1:]
	}

	r.pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(network, host, redis.DialConnectTimeout(time.Second), redis.DialReadTimeout(time.Second), redis.DialWriteTimeout(time.Second))
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

// Close the broker
func (r *RedisBroker) Close() error {
	return r.pool.Close()
}

// MigrateExpiredJobs is used for migrate expired jobs to ready queue
func (r *RedisBroker) MigrateExpiredJobs(queue string) {
	r.migrateExpiredJobs(fmt.Sprintf("%s:reserved", queue), queue)
	r.migrateExpiredJobs(fmt.Sprintf("%s:delayed", queue), queue)
}

func (r *RedisBroker) migrateExpiredJobs(from string, to string) {
	conn := r.pool.Get()
	defer conn.Close()
	if _, err := conn.Do("EVAL", MigrateJobsScript, 2, from, to, time.Now().Unix()); err != nil {
		logrus.Error(err)
		return
	}
	logrus.Debugf("migrate expired jobs from %s to %s", from, to)
}

// Reserve out a job [id] from broker with its life time
// if the reserved job is out of time(second)
// poller will kick it back in to ready queue
// if time out is 0, it means the job will be delete directly
func (r *RedisBroker) Reserve(queue string, timeout int64) (id string, err error) {
	conn := r.pool.Get()
	defer conn.Close()
	// reserve next job
	if timeout > 0 {
		id, err = redis.String(conn.Do("EVAL", ReserveScript, 2, queue, fmt.Sprintf("%s:reserved", queue), delayAt(timeout)))
	} else {
		id, err = redis.String(conn.Do("LPOP", queue))
	}
	if err != nil {
		return "", err
	}
	logrus.Debugf("pop job {%s} from queue {%s}", id, queue)
	return id, nil
}

// Delete the reserved job [id] from broker
// most of the time it means the job has been done successfully
func (r *RedisBroker) Delete(queue, id string) error {
	conn := r.pool.Get()
	defer conn.Close()
	if _, err := conn.Do("ZREM", fmt.Sprintf("%s:reserved", queue), id); err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Infof("delete reserved job {%s} from queue {%s}", id, queue)
	return nil
}

// Release is used for release the reserved job and push it back in to ready queue withe a delay(second) time
// this function maybe used for cron jobs
func (r *RedisBroker) Release(queue, id string, delay int64) error {
	conn := r.pool.Get()
	defer conn.Close()
	if _, err := conn.Do("EVAL", ReleaseScript, 2, fmt.Sprintf("%s:delayed", queue), fmt.Sprintf("%s:reserved", queue), id, delayAt(delay)); err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Infof("release job {%s} with delay {%d} to queue {%s}", id, delay, queue)
	return nil
}

// Push a job in to the queue
func (r *RedisBroker) Push(queue, id string) error {
	conn := r.pool.Get()
	defer conn.Close()
	if _, err := conn.Do("RPUSH", queue, id); err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Infof("push job {%s} to queue {%s}", id, queue)
	return nil
}

// Later is used for push a job in to the queue with a delay(second) time
// the job should be handled in the future time
func (r *RedisBroker) Later(queue, id string, delay int64) error {
	conn := r.pool.Get()
	defer conn.Close()
	if _, err := conn.Do("ZADD", fmt.Sprintf("%s:delayed", queue), delayAt(delay), id); err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Infof("push job {%s} with delay {%d} to queue {%s}", id, delay, queue)
	return nil
}

func delayAt(delay int64) int64 {
	return time.Now().Unix() + delay
}

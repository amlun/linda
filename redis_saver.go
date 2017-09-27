package linda

import (
	"github.com/garyburd/redigo/redis"
	neturl "net/url"
	"time"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"fmt"
)

const (
	// Job Info Prefix
	JobInfoPrefix = "linda:job:%s:info"
)

type RedisSaver struct {
	redisURL *neturl.URL
	pool     *redis.Pool
}

// Connect saver backend with url
func (r *RedisSaver) Connect(url *neturl.URL) error {
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
			c, err := redis.Dial(network, host)
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

// Close the saver
func (r *RedisSaver) Close() error {
	return r.pool.Close()
}

// Put the job to saver
func (r *RedisSaver) Put(job *Job) error {
	conn := r.pool.Get()
	defer conn.Close()
	bytes, err := json.Marshal(job)
	if err != nil {
		logrus.Error(err)
		return err
	}
	// job info
	key := fmt.Sprintf(JobInfoPrefix, job.ID)
	if _, err = conn.Do("SET", key, bytes); err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Infof("put job {%s}", job)
	return nil
}

// Get the job from saver
func (r *RedisSaver) Get(id string) (*Job, error) {
	conn := r.pool.Get()
	defer conn.Close()
	key := fmt.Sprintf(JobInfoPrefix, id)
	bytes, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	var job *Job
	if err := json.Unmarshal(bytes, &job); err != nil {
		logrus.Error(err)
		return nil, err
	}
	return job, nil
}

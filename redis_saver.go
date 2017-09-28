package linda

import (
	"github.com/garyburd/redigo/redis"
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
	url  string
	pool *redis.Pool
}

// Connect saver backend with url
func (r *RedisSaver) Connect(rawUrl string, timeout time.Duration) error {
	r.url = rawUrl
	r.pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(rawUrl, redis.DialConnectTimeout(timeout))
			if err != nil {
				return nil, err
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

// Delete the job from saver
func (r *RedisSaver) Delete(id string) error {
	conn := r.pool.Get()
	defer conn.Close()
	// job info
	key := fmt.Sprintf(JobInfoPrefix, id)
	if _, err := conn.Do("DEL", key); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

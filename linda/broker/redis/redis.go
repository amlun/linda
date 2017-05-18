package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/amlun/linda/linda/broker"
	"github.com/amlun/linda/linda/core"
	"github.com/garyburd/redigo/redis"
	neturl "net/url"
	"time"
)

type Broker struct {
	redisURL *neturl.URL
	pool     *redis.Pool
}

// connect redis with urlString redis://127.0.0.1:6379/0
func (b *Broker) Connect(url *neturl.URL) error {
	b.redisURL = url
	// redis connection pool
	b.pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", b.redisURL.Host)
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

func (b *Broker) Close() error {
	return b.pool.Close()
}

func (b *Broker) PushJob(job *core.Job) error {
	bytes, err := json.Marshal(job)
	if err != nil {
		return err
	}
	con := b.pool.Get()
	defer con.Close()
	if job.Queue == "" {
		return errors.New("queue can not be empty")
	}
	if _, err = con.Do("RPUSH", fmt.Sprintf("linda:queues:%s", job.Queue), bytes); err != nil {
		return err
	}
	//if job.Delay == 0 {
	//	job.RunTime = time.Now()
	//	con.Send("RPUSH", fmt.Sprintf("queues:%s", queue), bytes)
	//} else {
	//	job.RunTime = time.Now().Add(job.Delay)
	//	con.Send("ZADD", fmt.Sprintf("queues:%s:delayed", queue), job.RunTime.Unix(), bytes)
	//}
	return nil
}

func (b *Broker) GetJob(queue string) (*core.Job, error) {
	var job *core.Job
	con := b.pool.Get()
	defer con.Close()
	reply, err := redis.Bytes(con.Do("LPOP", fmt.Sprintf("linda:queues:%s", queue)))
	if err != nil {
		return nil, err
	}
	if reply == nil {
		return nil, errors.New("no job")
	}
	err = json.Unmarshal(reply, &job)
	if err != nil {
		return nil, err
	}
	return job, nil
}

//func (b *Broker) PushTask(task *core.Task) error {
//	bytes, err := json.Marshal(task)
//	if err != nil {
//		return err
//	}
//	con := b.pool.Get()
//	defer con.Close()
//	if _, err = con.Do("SADD", "linda:tasks", bytes); err != nil {
//		return err
//	}
//	return nil
//}

func (b *Broker) GetTask() (*core.Task, error) {
	var task *core.Task
	con := b.pool.Get()
	defer con.Close()
	reply, err := redis.Bytes(con.Do("SPOP", "linda:tasks"))
	if err != nil {
		return nil, err
	}
	if reply == nil {
		return nil, errors.New("no task")
	}
	err = json.Unmarshal(reply, &task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (b *Broker) Length(queue string) int {
	con := b.pool.Get()
	defer con.Close()
	length, err := redis.Int(con.Do("LLEN", fmt.Sprintf("linda:queues:%s", queue)))
	if err != nil {
		return 0
	}
	return length
}

// register redis broker when module init
func init() {
	broker.Register("redis", &Broker{})
}

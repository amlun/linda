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

func (b *Broker) PushJob(queue string, job *core.Job) error {
	job.RunTime = time.Now()
	bytes, err := json.Marshal(job)
	if err != nil {
		return err
	}
	con := b.pool.Get()
	defer con.Close()
	if queue == "" {
		return errors.New("queue can not be empty")
	}
	_, err = con.Do("RPUSH", fmt.Sprintf("queues:%s", queue), bytes)
	if err != nil {
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

func (b *Broker) GetJob(queue string, job *core.Job) error {
	con := b.pool.Get()
	defer con.Close()
	reply, err := redis.Bytes(con.Do("LPOP", fmt.Sprintf("queues:%s", queue)))
	if err != nil {
		return err
	}
	err = json.Unmarshal(reply, job)
	if err != nil {
		return err
	}
	return nil

}

func (b *Broker) Length(queue string) int {
	con := b.pool.Get()
	defer con.Close()
	length, err := redis.Int(con.Do("LLEN", fmt.Sprintf("queues:%s", queue)))
	if err != nil {
		return 0
	}
	return length
}

// register redis broker when module init
func init() {
	broker.Register("redis", &Broker{})
}

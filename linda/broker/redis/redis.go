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

func (b *Broker) QueueMonitors() []core.QueueStatus {
	con := b.pool.Get()
	defer con.Close()
	var queueStatus core.QueueStatus
	queues, _ := redis.Strings(con.Do("SMEMBERS", "queues"))
	queueStatusList := make([]core.QueueStatus, len(queues))
	for i, queue := range queues {
		length, _ := redis.Int(con.Do("LLEN", fmt.Sprintf("queues:%s", queue)))
		queueStatus.Queue = queue
		queueStatus.Length = length
		queueStatusList[i] = queueStatus
	}
	return queueStatusList
}

func (b *Broker) PushJob(job *core.Job) error {
	job.RunTime = time.Now()
	bytes, err := json.Marshal(job)
	if err != nil {
		return err
	}
	con := b.pool.Get()
	defer con.Close()
	queue := job.Func
	if queue == "" {
		return errors.New("queue can not be empty")
	}
	// All queues for monitor
	con.Send("SADD", "queues", queue)
	con.Send("RPUSH", fmt.Sprintf("queues:%s", queue), bytes)
	//if job.Delay == 0 {
	//	job.RunTime = time.Now()
	//	con.Send("RPUSH", fmt.Sprintf("queues:%s", queue), bytes)
	//} else {
	//	job.RunTime = time.Now().Add(job.Delay)
	//	con.Send("ZADD", fmt.Sprintf("queues:%s:delayed", queue), job.RunTime.Unix(), bytes)
	//}
	con.Flush()
	return nil
}

func (b *Broker) GetJob(queue string, job *core.Job) error {
	con := b.pool.Get()
	defer con.Close()
	if queue == "" {
		return errors.New("queue can not be empty")
	}
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

// register redis broker when module init
func init() {
	broker.Register("redis", &Broker{})
}

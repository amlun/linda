package redis

import (
	"github.com/amlun/linda/linda/smarter"
	"github.com/garyburd/redigo/redis"
	neturl "net/url"
	"time"
)

type Smarter struct {
	redisURL *neturl.URL
	pool     *redis.Pool
}

// connect redis with urlString redis://127.0.0.1:6379/0
func (s *Smarter) Connect(url *neturl.URL) error {
	s.redisURL = url
	// redis connection pool
	s.pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", s.redisURL.Host)
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

func (s *Smarter) Close() error {
	return s.pool.Close()
}

func (s *Smarter) PushTask(taskId string) error {
	con := s.pool.Get()
	defer con.Close()
	if _, err := con.Do("SADD", "linda:tasks", taskId); err != nil {
		return err
	}
	return nil
}

func (s *Smarter) GetTask() (taskId string, err error) {
	con := s.pool.Get()
	defer con.Close()
	taskId, err = redis.String(con.Do("SPOP", "linda:tasks"))
	if redis.ErrNil == err {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return taskId, err
}

// register redis broker when module init
func init() {
	smarter.Register("redis", &Smarter{})
}

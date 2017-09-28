package linda

import (
	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func setupRedisSaverTestCase() {
	logrus.SetLevel(logrus.DebugLevel)
}

func TestRedisSaver_Connect(t *testing.T) {
	setupRedisSaverTestCase()
	saver = &RedisSaver{}
	if err := saver.Connect("redis://localhost:6379/", time.Second); err != nil {
		t.Error(err)
		return
	}

	t.Log(saver)

}

func TestRedisSaver_Close(t *testing.T) {
	TestRedisSaver_Connect(t)
	if err := saver.Close(); err != nil {
		t.Error(err)
		return
	}
	t.Log("TestRedisSaver_Close")
}

func TestRedisSaver_Get(t *testing.T) {
	TestRedisSaver_Connect(t)
	job, err := saver.Get("1")
	if err != nil && err != redis.ErrNil {
		t.Error(err)
		return
	}
	t.Log(job)
}

func TestRedisSaver_Put(t *testing.T) {
	TestRedisSaver_Connect(t)
	job := &Job{
		ID:        "1",
		Queue:     "test",
		Period:    60,
		Retry:     3,
		CreatedAt: time.Now(),
	}
	err := saver.Put(job)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(job)
}

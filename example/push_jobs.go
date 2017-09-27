package main

import (
	"github.com/amlun/linda"
	"github.com/sirupsen/logrus"
	"time"
)

func main() {
	var err error
	var b linda.Broker
	var s linda.Saver
	// broker
	if b, err = linda.NewBroker("redis://localhost:6379/"); err != nil {
		logrus.Error(err)
		return
	}
	// saver
	if s, err = linda.NewSaver("redis://localhost:6379/"); err != nil {
		logrus.Error(err)
		return
	}
	// job
	var jobID = "1"
	var queue = "test"
	var job = &linda.Job{
		ID:        jobID,
		Queue:     queue,
		Period:    60,
		Retry:     3,
		CreatedAt: time.Now(),
		Payload: linda.Payload{
			Class: "printArgs",
			Args:  []interface{}{"a", "b", "c"},
		},
	}
	// save job
	if err = s.Put(job); err != nil {
		logrus.Error(err)
		return
	}
	// push to broker
	if err = b.Push(queue, jobID); err != nil {
		logrus.Error(err)
		return
	}
}

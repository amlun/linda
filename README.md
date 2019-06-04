# Linda

[![Build Status](https://travis-ci.org/amlun/linda.svg?branch=master)](https://travis-ci.org/amlun/linda)
[![GoDoc](https://godoc.org/github.com/amlun/linda?status.svg)](https://godoc.org/github.com/amlun/linda)
[![GoReport](https://goreportcard.com/badge/github.com/amlun/linda)](https://goreportcard.com/report/github.com/amlun/linda)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://raw.githubusercontent.com/amlun/linda/master/LICENSE)

Linda is a background manager to poll jobs from broker and dispatch them to multi workers.

Linda Broker provides a unified API across different broker (queue) services.

Linda Saver provides a unified API across different saver (db) services.

Brokers allow you to defer the processing of a time consuming task.

When job done, use Release func to release the job with a delay (seconds), you can implement a cron job service.

The real period is job.Period + Interval

Inspiration comes from [beanstalkd](https://github.com/kr/beanstalkd) and [goworker](https://github.com/benmanns/goworker) 

## Installation

To install Linda, use
```sh
go get github.com/amlun/linda
```
to get the main package, and then use [glide](https://glide.sh/)
```sh
glide install
```
to install the dependencies

## Getting Started

### Terminology

* Broker
> message transport [MQ]

* Saver
> job info storage [Database]

* poller
> poll job from the broker and send to local job channels

> poller also migrate the expire jobs

* worker
> worker is the main process to work the job

### Worker Type

```
type workerFunc func(...interface{}) error
```

### Register Worker
```
linda.RegisterWorkers("MyClass", myFunc)
```

### Broker Interface
```
type Broker interface {
	Connect(url *neturl.URL) error
	Close() error
	MigrateExpiredJobs(queue string)
	Reserve(queue string, timeout int64) (string, error)
	Delete(queue, id string) error
	Release(queue, id string, delay int64) error
	Push(queue, id string) error
	Later(queue, id string, delay int64) error
}
```

### Saver Interface
```
type Saver interface {
	Connect(url *neturl.URL) error
	Close() error
	Put(job *Job) error
	Get(id string) (*Job, error)
}
```

### Examples

Add jobs to saver and push them to broker

```sh
go run example/push_jobs/main.go
```

example/push_jobs/main.go

```go
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

```

Worker run to consume the job

```sh
go run example/print_args/main.go
```

example/print_args/main.go

```go
package main

import (
	"fmt"
	"github.com/amlun/linda"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	linda.RegisterWorkers("printArgs", PrintArgs)
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	// broker
	b, _ := linda.NewBroker("redis://localhost:6379/")
	// saver
	s, _ := linda.NewSaver("redis://localhost:6379/")
	// config
	c := linda.Config{
		Queue:     "test",
		Timeout:   60,
		Interval:  time.Second,
		WorkerNum: 4,
	}
	quit := signals()
	linda.Init(c, b, s)
	go func() {
		defer func() {
			linda.Quit()
		}()
		<-quit
	}()

	if err := linda.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}

func PrintArgs(args ...interface{}) error {
	fmt.Println(args)
	return nil
}

// Signal Handling
func signals() <-chan bool {
	quit := make(chan bool)
	go func() {
		signals := make(chan os.Signal)
		defer close(signals)
		signal.Notify(signals, syscall.SIGQUIT, syscall.SIGTERM, os.Interrupt)
		defer signalStop(signals)
		<-signals
		quit <- true
	}()
	return quit
}

// Stops signals channel.
func signalStop(c chan<- os.Signal) {
	signal.Stop(c)
}

```

## Features

### Broker List

 - [x] Redis
 - [ ] beanstalkd
 - [ ] RabbitMQ
 
## Design

### System Design

![system-design](https://rawgit.com/amlun/linda/master/images/linda.png)

### Job State
```
   later                                release
  ----------------> [DELAYED] <------------.
                        |                   |
                   kick | (time passes)     |
                        |                   |
   push                 v     reserve       |       delete
  -----------------> [READY] ---------> [RESERVED] --------> *poof*
                        ^                   |
                         \                  |
                          `-----------------'
                           kick (time out)
 
```
## Thanks

* [redigo](https://github.com/garyburd/redigo)

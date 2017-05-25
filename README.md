# Linda

[![Build Status](https://travis-ci.org/amlun/linda.png?branch=master)](https://travis-ci.org/amlun/linda)
[![GoDoc](https://godoc.org/github.com/amlun/linda?status.svg)](https://godoc.org/github.com/amlun/linda)

Linda is a background manager to poll jobs from broker and dispatch them to multi workers.

Linda Broker provides a unified API across different broker (queue) services.

Brokers allow you to defer the processing of a time consuming task.

Use ReleaseWithDelay func, you can implement a cron job service.

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

* poller
> poll job from the broker and send to local job channels

> poller also migrate the expire jobs

* worker
> worker is the main process to work the job

### Worker Type

```
func(job *Job) error
```

### Register Worker
```
linda.Register("MyClass", myFunc)
```

### Broker Interface
```go
package linda

import (
	neturl "net/url"
)

type Broker interface {
	Connect(url *neturl.URL) error
	Close() error
	MigrateExpiredJobs(queue string)
	Pop(queue string, timeout int64) (*Job, error)
	Delete(queue string, job *Job) error
	Release(queue string, job *Job) error
	ReleaseWithDelay(queue string, job *Job, delay int64) error
	Push(job *Job, queue string) error
	Later(delay int64, job *Job, queue string) error
}
```

### Examples


push a job to queue

```go
package main

import (
	"fmt"
	"github.com/amlun/linda"
)

func main() {
	broker, err := linda.NewBroker("redis://localhost:6379/")
	if err != nil {
		fmt.Println(err)
		return
	}
	queue := "scheduler"
	job := &linda.Job{
		Queue: queue,
		Payload: linda.Payload{
			Class: "DispatcherSeed",
			Args:  []interface{}{"seed_url_md5"},
		},
	}

	if err := broker.Push(job, queue); err != nil {
		fmt.Println("Error:", err)
		return
	}
}

```

Worker run to consume the job
```go
package main

import (
	"fmt"
	"github.com/amlun/linda"
)

func init() {
	settings := linda.Settings{
		Queue:         "scheduler",
		Connection:      "redis://localhost:6379/",
		Timeout:       60,
		IntervalFloat: 5.0,
		Concurrency:   1,
	}
	linda.SetSettings(settings)
	linda.RegisterWorkers("DispatcherSeed", DispatcherSeed)
}

func main() {
	if err := linda.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}

func DispatcherSeed(job *linda.Job) error {
	broker := linda.GetBroker()
	// get seed info
	// do seed job
	// release job with delay (like a cron job)
	broker.ReleaseWithDelay("scheduler", job, 60)
	return nil
}
```

## Features

### Broker List

 - [x] Redis
 - [ ] beanstalkd
 - [ ] NSQ
 - [ ] Kafka
 - [ ] RabbitMQ
 
## Design

### System Design

![system-design](https://rawgit.com/amlun/linda/master/images/linda.png)

### Job State
```
   later                       release with delay
  ----------------> [DELAYED] <------------.
                        |                   |
                migrate | (time passes)     |
                        |                   |
   push                 v     reserve       |       delete
  -----------------> [READY] ---------> [RESERVED] --------> *poof*
                        ^                |
                         \    release    |
                          `--------------'
                           migrate (time out)
 
```
## Thanks

* [redigo](https://github.com/garyburd/redigo)
* [goworker](https://github.com/benmanns/goworker)
* [laravel/queue](https://github.com/laravel/framework)
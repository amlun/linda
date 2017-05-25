# Linda

[![Build Status](https://travis-ci.org/amlun/linda.svg?branch=master)](https://travis-ci.org/amlun/linda)
[![GoDoc](https://godoc.org/github.com/amlun/linda?status.svg)](https://godoc.org/github.com/amlun/linda)
[![GoReport](https://goreportcard.com/badge/github.com/amlun/linda)](https://goreportcard.com/report/github.com/amlun/linda)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://raw.githubusercontent.com/amlun/linda/master/LICENSE)

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


use redis-cli push jobs to queue

```
RPUSH print "{\"queue\":\"print\",\"Payload\":{\"class\":\"printArgs\",\"args\":[\"a\",\"b\",\"c\"]}}"
RPUSH print "{\"queue\":\"print\",\"Payload\":{\"class\":\"printArgs\",\"args\":[\"A\",\"B\",\"C\"]}}"
RPUSH print "{\"queue\":\"print\",\"Payload\":{\"class\":\"printArgs\",\"args\":[1,2,3,4,5,6,7]}}"
```

Worker run to consume the job
```sh
go run example/print_args -queue=print -connection=redis://localhost:6379/
```

example/print_args.go

```go
package main

import (
	"fmt"
	"github.com/amlun/linda"
)

func init() {
	linda.RegisterWorkers("PrintArgs", PrintArgs)
}

func main() {
	if err := linda.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}

func PrintArgs(job *linda.Job) error {
	fmt.Println(job.Payload.Args)
	linda.GetBroker().Delete(job.Queue, job)
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
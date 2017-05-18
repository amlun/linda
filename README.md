# Linda

Linda is a simple dispatcher and scheduler system.

Linda has several tools/apis to manage tasks and jobs.

It allows you to save tasks to saver[DB] and push jobs to broker[MQ].

## Installation

To install Linda, use 

`go get github.com/amlun/linda`

to install the package, and then use [glide](https://glide.sh/)

`glide install`

to install the dependency packages

## Getting Started

### Terminology

* Broker
> message transport [MQ]

* Smarter
> save scheduled tasks and schedule them

* Saver
> backend to store all things

* Queues
> receive job and send to exchanges

* Periods
> period of the task

* Task
> template of job

* Job
> job is a callable class/function with args 


### Simple Usage

Edit the apps/server.go, apps/scheduler.go, modify the config with your own urlString

```go
package main

import (
	"github.com/amlun/linda/linda"
	"github.com/amlun/linda/modules/server"
)

func main() {
	var config = linda.Config{
		BrokerURL: "redis://127.0.0.1:6379",
		SaverURL:  "cassandra://cassandra:cassandra@127.0.0.1:9042/linda",
		SmarterURL: "redis://127.0.0.1:6379",
	}
	l := linda.NewLinda(&config)
	defer l.Close()
	server.Start(l)
}
```

```go
package main

import (
	"github.com/amlun/linda/linda"
	"github.com/amlun/linda/modules/scheduler"
)

func main() {
	var config = linda.Config{
		BrokerURL: "redis://127.0.0.1:6379",
		SaverURL:  "cassandra://cassandra:cassandra@127.0.0.1:9042/linda",
		SmarterURL: "redis://127.0.0.1:6379",
	}
	l := linda.NewLinda(&config)
	defer l.Close()
	s := scheduler.New(l)
	s.Start()
}

```

And then use

`go run apps/server.go`

to start a http server and serve the apis.

Use

`go run apps/scheduler.go`

to start scheduler to schedule the periodic task.


### API Doc

 * GET /api/ping - Check the server if it is alive.
 * GET /api/tasks - List all tasks.
 * GET /api/queues - List all queues and its length.
 * GET /api/job - Get a job from queue, now it only implements a simple way to fetch a job.
 * POST /api/task - Post a task.
 
### API Usage 

#### Post A Task
HTTP method: `POST`

Host/port: `http://localhost:8081/api/task`

Request Parameters: `Func=test&Args=a&Args=b&Args=c&Period=100`

#### Get A Job
HTTP method: `GET`

Host/port: `http://localhost:8081/api/job`

Request Parameters: `queue=test`

### Job in Queue

```json
{
    "job_id": "9964015c-a96c-4e48-aad6-985ab8dc1888",
    "run_time": "2017-05-17T14:19:50.355512848+08:00",
    "delay": 0,
    "status": 0,
    "task_id": "57354766-93f6-43eb-8b86-a4bc2b547cf8",
    "period": 100,
    "func": "test",
    "args": [
        "a",
        "b",
        "c"
    ]
}
```
The only useful field is `args`, because `func` is same with `queue` name.

Clients fetch jobs from `queue`, and then handle it with `args`.

And then update the job status. @TODO

Job_Status:

```
	PENDING = iota
	STARTED 
	RETRY   
	SUCCESS 
	FAILURE 
```

## Features

### Linda Dispatcher

 - [x] Simple dispatcher with func as the queue
 - [ ] Manage a queue pool and support priority queue
 
### Linda Scheduler

 - [x] Simple scheduler with multi periods, one period with one schedule worker
 - [ ] Distribute scheduler workers

### Broker List

 - [x] Redis
 - [ ] NSQ
 - [ ] Kafka
 - [ ] RabbitMQ

### Saver List

 - [x] Cassandra
 - [ ] Redis
 - [ ] Mysql
 
### Web UI

 - [ ] Task List & Manage
 - [ ] Job List & Manage
 - [ ] Queue List & Manage & Monitor
 - [ ] Periods Manage / Cron Jobs
 - [ ] Data Statistics
 
### Clients
 - [x] HTTP API
 - [ ] Go
 - [ ] Python
 
## Thanks

* [redigo](https://github.com/garyburd/redigo)
* [gocql](https://github.com/gocql/gocql)
* [scheduler](https://github.com/carlescere/scheduler)
* [Gin](https://github.com/gin-gonic/gin)
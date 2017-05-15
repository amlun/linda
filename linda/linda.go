package linda

import (
	"github.com/amlun/linda/linda/core"
	"github.com/twinj/uuid"
	"log"
	"time"
)

type Linda struct {
	config *Config
	dispatcher
}

func NewLinda(config *Config) *Linda {
	if config.BrokerURL == "" {
		config.BrokerURL = "redis://localhost:6379"
	}
	if config.SaverURL == "" {
		config.SaverURL = "cassandra://localhost:9042/linda"
	}
	l := &Linda{
		config: config,
		dispatcher: dispatcher{
			brokerURL: config.BrokerURL,
			saverURL:  config.SaverURL,
		},
	}
	if l.dispatcher.Init() != nil {
		panic("Linda dispatcher init failed")
	}
	return l
}

func (l *Linda) Close() {
	l.dispatcher.Close()
}

// schedule jobs with frequency
func (l *Linda) Schedule(frequency time.Duration) func() {
	return func() {
		var job core.Job
		tasks := make(chan core.Task)
		go func() {
			l.saver.GetTimingTask(frequency, tasks)
		}()
		for task := range tasks {
			l.saver.ScheduleTask(task.ID)
			job.ID = uuid.NewV4().String()
			job.TaskId = task.ID
			job.Func = task.Func
			job.Args = task.Args
			l.PushJob(job)
		}
		log.Printf("schedule the job with frequency [%s]", frequency)
	}
}

// schedule list
func (l *Linda) Schedules() []time.Duration {
	return l.saver.Frequencies()
}

// get all task queues and monitor
func (l *Linda) MonitorQueues() []core.QueueStatus {
	return l.broker.QueueMonitors()
}

// get task list
func (l *Linda) TaskList(taskList *core.TaskList) error {
	return l.saver.TaskList(taskList)
}

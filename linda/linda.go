package linda

import (
	"github.com/amlun/linda/linda/core"
	"github.com/twinj/uuid"
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
	if config.SmarterURL == "" {
		config.SmarterURL = "redis://localhost:6379"
	}
	l := &Linda{
		config: config,
		dispatcher: dispatcher{
			brokerURL:  config.BrokerURL,
			saverURL:   config.SaverURL,
			smarterURL: config.SmarterURL,
		},
	}
	if err := l.dispatcher.Init(); err != nil {
		panic(err)
	}
	return l
}

func (l *Linda) ScheduleTask(task *core.Task) func() {
	return func() {
		var job core.Job
		job.JobId = uuid.NewV4().String()
		job.TaskId = task.TaskId
		job.Queue = task.Queue
		job.Payload = task.Payload
		l.PushJob(job)
	}
}

func (l *Linda) Close() {
	l.dispatcher.Close()
}

// get task list
func (l *Linda) TaskList(taskList *core.TaskList) error {
	return l.saver.TaskList(taskList)
}

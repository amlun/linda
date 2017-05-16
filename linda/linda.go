package linda

import (
	"github.com/amlun/linda/linda/core"
	log "github.com/sirupsen/logrus"
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
	l := &Linda{
		config: config,
		dispatcher: dispatcher{
			brokerURL: config.BrokerURL,
			saverURL:  config.SaverURL,
		},
	}
	if l.dispatcher.Init() != nil {
		log.Panic("Linda dispatcher init failed")
	}
	return l
}

func (l *Linda) Close() {
	l.dispatcher.Close()
}

// schedule jobs with period
func (l *Linda) Schedule(period int) func() {
	return func() {
		var job core.Job
		tasks := make(chan core.Task)
		go func() {
			l.saver.GetPeriodicTask(period, tasks)
		}()
		for task := range tasks {
			l.saver.ScheduleTask(task.TaskId)
			job.JobId = uuid.NewV4().String()
			job.Task = task
			l.PushJob(job)
		}
		log.WithField("period", period).Info("schedule the job with period")
	}
}

// schedule list
func (l *Linda) Periods() []int {
	return l.saver.Periods()
}

// get all task queues and monitor
func (l *Linda) MonitorQueues() []core.QueueStatus {
	var queueStatus core.QueueStatus
	queues := l.saver.Queues()
	queueStatusList := make([]core.QueueStatus, len(queues))
	for i, queue := range queues {
		queueStatus.Queue = queue
		queueStatus.Length = l.broker.Length(queue)
		queueStatusList[i] = queueStatus
	}
	return queueStatusList
}

// get task list
func (l *Linda) TaskList(taskList *core.TaskList) error {
	return l.saver.TaskList(taskList)
}

package scheduler

import (
	"github.com/amlun/linda/linda"
	cron "github.com/carlescere/scheduler"
	"time"
)

var Linda *linda.Linda

type scheduler struct {
}

func New() *scheduler {
	return &scheduler{}
}

func (s *scheduler) Start(linda *linda.Linda) error {
	// init linda instance
	Linda = linda
	// register signals
	quit := signals()
	poller, err := newPoller()
	if err != nil {
		return err
	}
	// poll task ids
	taskIds := poller.poll(time.Second*2, quit)
	for taskId := range taskIds {
		task, err := Linda.GetTask(taskId)
		if err != nil {
			return err
		}
		cron.Every(task.Period).Seconds().Run(Linda.ScheduleTask(task))
	}
	return nil

}

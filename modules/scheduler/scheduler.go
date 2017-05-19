package scheduler

import (
	"github.com/amlun/linda/linda"
	"time"
)

var (
	Linda *linda.Linda
	log   = linda.Logger
)

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
	// worker run cron jobs
	worker, err := newWorker()
	if err != nil {
		return err
	}
	worker.work(taskIds)
	return nil

}

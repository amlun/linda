package scheduler

import (
	"github.com/amlun/linda/linda"
	cron "github.com/carlescere/scheduler"
)

type scheduler struct {
	linda    *linda.Linda
	isMaster bool
}

func New(linda *linda.Linda) *scheduler {
	return &scheduler{
		linda: linda,
	}
}

// will support distribute deploy
// use redis or zookeeper to lock ,one master with multi slave
func (s *scheduler) Start() {
	// register()
	quit := signals()
	list := s.linda.Periods()
	for _, period := range list {
		cron.Every(period).Seconds().Run(s.linda.Schedule(period))
	}
	<-quit
}

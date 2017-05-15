package monitor

import (
	"fmt"
	"github.com/amlun/linda/linda"
	cron "github.com/carlescere/scheduler"
	"runtime"
)

type monitor struct {
	linda *linda.Linda
}

func New(linda *linda.Linda) *monitor {
	return &monitor{
		linda: linda,
	}
}

//
//
func (m *monitor) Start() {
	job := func() {
		data := m.linda.MonitorQueues()
		fmt.Println(data)
	}
	cron.Every(1).Seconds().Run(job)
	runtime.Goexit()
}

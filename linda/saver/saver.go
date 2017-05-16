package saver

import (
	"fmt"
	"github.com/amlun/linda/linda/core"
	neturl "net/url"
)

type Saver interface {
	Connect(url *neturl.URL) error
	Close() error
	PublishTask(t *core.Task) error
	PublishJob(t *core.Job) error
	Periods() []int
	Queues() []string
	GetPeriodicTask(period int, tasks chan core.Task)
	ScheduleTask(id string) error
	TaskList(taskList *core.TaskList) error
}

// registered savers
var saverRegistery = make(map[string]Saver)

// Register saver with its scheme
// For example: mysql://127.0.0.1:3306
func Register(scheme string, s Saver) {
	saverRegistery[scheme] = s
}

func NewSaver(urlString string) (Saver, error) {
	url, err := neturl.Parse(urlString)
	if err != nil {
		return nil, err
	}
	// get scheme from uri
	var scheme = url.Scheme
	if s, ok := saverRegistery[scheme]; ok {
		err := s.Connect(url)
		if err != nil {
			return nil, err
		}
		return s, nil
	}
	return nil, fmt.Errorf("Unknow saver scheme [%s]", scheme)
}

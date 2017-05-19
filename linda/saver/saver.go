package saver

import (
	"fmt"
	"github.com/amlun/linda/linda/core"
	neturl "net/url"
)

type Saver interface {
	Connect(url *neturl.URL) error
	Close() error
	SaveTask(t *core.Task) error
	SaveJob(t *core.Job) error
	GetTask(taskId string) (*core.Task, error)
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

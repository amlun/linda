package smarter

import (
	"fmt"
	neturl "net/url"
)

type Smarter interface {
	Connect(url *neturl.URL) error
	Close() error
	PushTask(taskId string) error
	GetTask() (taskId string, err error)
}

// registered smarters
var smarterRegistery = make(map[string]Smarter)

// Register smarter with its scheme
func Register(scheme string, b Smarter) {
	smarterRegistery[scheme] = b
}

func NewSmarter(urlString string) (Smarter, error) {
	// get scheme from uri
	url, err := neturl.Parse(urlString)
	if err != nil {
		return nil, err
	}
	scheme := url.Scheme
	if s, ok := smarterRegistery[scheme]; ok {
		err := s.Connect(url)
		if err != nil {
			return nil, err
		}
		return s, nil
	}

	return nil, fmt.Errorf("Unknow smarter scheme [%s]", scheme)
}

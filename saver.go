package linda

import (
	"errors"
	neturl "net/url"
	"github.com/sirupsen/logrus"
)

var (
	UnknownSaver = errors.New("unknown saver scheme")
)

// Saver is job saved on [Storage]
// it provides a unified API, support multi drivers
type Saver interface {
	Connect(url *neturl.URL) error
	Close() error
	Put(job *Job) error
	Get(id string) (*Job, error)
}

var saverMaps = make(map[string]Saver)

// RegisterSaver is used to register savers with scheme name
// You can use your own saver driver
func RegisterSaver(scheme string, saver Saver) {
	logrus.Debugf("register saver [%s]", scheme)
	if saver == nil {
		panic("Register saver is nil")
	}
	saverMaps[scheme] = saver
}

// NewSaver will get an instance of saver with url string
// if there is no matched scheme, return error
// now saver only support redis
func NewSaver(urlString string) (Saver, error) {
	url, err := neturl.Parse(urlString)
	if err != nil {
		return nil, err
	}
	scheme := url.Scheme
	if s, ok := saverMaps[scheme]; ok {
		err := s.Connect(url)
		if err != nil {
			return nil, err
		}
		return s, nil
	}
	return nil, UnknownSaver
}

func init() {
	RegisterSaver("redis", &RedisSaver{})
}

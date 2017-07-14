package linda

import (
	"fmt"
	neturl "net/url"
)

// Broker is message transport[MQ]
// it provides a unified API, support multi drivers
type Broker interface {
	Connect(url *neturl.URL) error
	Close() error
	MigrateExpiredJobs(queue string)
	Reserve(queue string, timeout int64) (*Job, error)
	Delete(queue string, job *Job) error
	Release(queue string, job *Job, delay int64) error
	Push(job *Job, queue string) error
	Later(delay int64, job *Job, queue string) error
}

var brokerMaps = make(map[string]Broker)

// RegisterBroker is used to register brokers with scheme name
// You can use your own broker driver
func RegisterBroker(scheme string, broker Broker) {
	if broker == nil {
		panic("Register broker is nil")
	}
	brokerMaps[scheme] = broker
}

// NewBroker will get an instance of broker with url string
// if there is no matched scheme, return error
// now broker only support redis
func NewBroker(urlString string) (Broker, error) {
	url, err := neturl.Parse(urlString)
	if err != nil {
		return nil, err
	}
	scheme := url.Scheme
	if b, ok := brokerMaps[scheme]; ok {
		err := b.Connect(url)
		if err != nil {
			return nil, err
		}
		return b, nil
	}
	return nil, fmt.Errorf("Unknow broker scheme [%s]", scheme)
}

func init() {
	RegisterBroker("redis", &RedisBroker{})
}

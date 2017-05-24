package linda

import (
	"fmt"
	neturl "net/url"
)

type Broker interface {
	Connect(url *neturl.URL) error
	Close() error
	MigrateExpiredJobs(queue string)
	Pop(queue string, ack bool, timeout int64) (*Job, error)
	DeleteReserved(queue string, job *Job) error
	DeleteAndRelease(queue string, job *Job, delay int64) error
	Push(job *Job, queue string) error
	Later(delay int64, job *Job, queue string) error
}

var brokerRegistery = make(map[string]Broker)

func RegisterBroker(scheme string, broker Broker) {
	brokerRegistery[scheme] = broker
}

func NewBroker(urlString string) (Broker, error) {
	url, err := neturl.Parse(urlString)
	if err != nil {
		return nil, err
	}
	scheme := url.Scheme
	if b, ok := brokerRegistery[scheme]; ok {
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

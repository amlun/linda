package broker

import (
	"fmt"
	"github.com/amlun/linda/linda/core"
	neturl "net/url"
)

type Broker interface {
	Connect(url *neturl.URL) error
	Close() error
	PushJob(job *core.Job) error
	GetJob(queue string, job *core.Job) error
	Length(queue string) int
}

// registered brokers
var brokerRegistery = make(map[string]Broker)

// Register broker with its scheme
func Register(scheme string, b Broker) {
	brokerRegistery[scheme] = b
}

func NewBroker(urlString string) (Broker, error) {
	// get scheme from uri
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

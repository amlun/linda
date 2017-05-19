package linda

import (
	"github.com/amlun/linda/linda/broker"
	_ "github.com/amlun/linda/linda/broker/redis"
	"github.com/amlun/linda/linda/saver"
	_ "github.com/amlun/linda/linda/saver/cassandra"
	"github.com/amlun/linda/linda/smarter"
	_ "github.com/amlun/linda/linda/smarter/redis"
)

type Linda struct {
	config  *Config
	broker  broker.Broker
	saver   saver.Saver
	smarter smarter.Smarter
}

func NewLinda(config *Config) *Linda {
	if config.BrokerURL == "" {
		config.BrokerURL = "redis://localhost:6379"
	}
	if config.SaverURL == "" {
		config.SaverURL = "cassandra://localhost:9042/linda"
	}
	if config.SmarterURL == "" {
		config.SmarterURL = "redis://localhost:6379"
	}
	l := &Linda{
		config: config,
	}
	if err := l.init(); err != nil {
		panic(err)
	}
	return l
}

func (l *Linda) init() error {
	var err error
	if l.broker, err = broker.NewBroker(l.config.BrokerURL); err != nil {
		Logger.Error(err)
		return err
	}
	if l.saver, err = saver.NewSaver(l.config.SaverURL); err != nil {
		Logger.Error(err)
		return err
	}
	if l.smarter, err = smarter.NewSmarter(l.config.SmarterURL); err != nil {
		Logger.Error(err)
		return err
	}
	return nil
}

func (l *Linda) Close() {
	l.broker.Close()
	l.saver.Close()
	l.smarter.Close()
}

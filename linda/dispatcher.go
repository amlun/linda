package linda

import (
	"github.com/amlun/linda/linda/broker"
	_ "github.com/amlun/linda/linda/broker/redis"
	"github.com/amlun/linda/linda/core"
	"github.com/amlun/linda/linda/saver"
	_ "github.com/amlun/linda/linda/saver/cassandra"
	log "github.com/sirupsen/logrus"
)

type dispatcher struct {
	brokerURL string
	broker    broker.Broker
	saverURL  string
	saver     saver.Saver
}

func (d *dispatcher) Init() error {
	b, err := broker.NewBroker(d.brokerURL)
	if err != nil {
		log.Error(err)
		return err
	}
	d.broker = b
	s, err := saver.NewSaver(d.saverURL)
	if err != nil {
		log.Error(err)
		return err
	}
	d.saver = s
	return nil
}

func (d *dispatcher) Close() {
	d.broker.Close()
	d.saver.Close()
}

// push a [frequency] task to saver
func (d *dispatcher) PushTask(task core.Task) error {
	err := d.saver.PublishTask(&task)
	if err != nil {
		log.Error(err)
		return err
	}
	log.WithField("task", task).Info("push task to saver")
	if task.Frequency != 0 {
		// save task frequency for scheduler
		err = d.saver.Frequency(task.Frequency)
		if err != nil {
			return err
		}
		log.WithField("task", task).Info("save task frequency for scheduler")
	}
	return nil
}

// push a job to broker and save it
func (d *dispatcher) PushJob(job core.Job) error {
	err := d.broker.PushJob(&job)
	if err != nil {
		log.Error(err)
		return err
	}
	d.saver.PublishJob(&job)
	log.WithField("job", job).Info("push job to broker and saver")
	return nil
}

// get a job and delete it from the queue
func (d *dispatcher) GetJob(queue string) core.Job {
	var job core.Job
	err := d.broker.GetJob(queue, &job)
	if err != nil {
		log.Error(err)
	}
	return job

}

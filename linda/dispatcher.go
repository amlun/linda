package linda

import (
	"github.com/amlun/linda/linda/broker"
	_ "github.com/amlun/linda/linda/broker/redis"
	"github.com/amlun/linda/linda/core"
	"github.com/amlun/linda/linda/saver"
	_ "github.com/amlun/linda/linda/saver/cassandra"
	"github.com/amlun/linda/linda/smarter"
	_ "github.com/amlun/linda/linda/smarter/redis"
)

type dispatcher struct {
	brokerURL  string
	broker     broker.Broker
	saverURL   string
	saver      saver.Saver
	smarterURL string
	smarter    smarter.Smarter
}

func (d *dispatcher) Init() error {
	var err error
	if d.broker, err = broker.NewBroker(d.brokerURL); err != nil {
		Logger.Error(err)
		return err
	}
	if d.saver, err = saver.NewSaver(d.saverURL); err != nil {
		Logger.Error(err)
		return err
	}
	if d.smarter, err = smarter.NewSmarter(d.smarterURL); err != nil {
		Logger.Error(err)
		return err
	}
	return nil
}

func (d *dispatcher) Close() {
	d.broker.Close()
	d.saver.Close()
}

// push a [period] task to saver
// if this is a scheduled task, then push to smarter
func (d *dispatcher) PushTask(task core.Task) error {
	log := Logger.WithField("action", "PushTask").WithField("task", task)
	if err := d.saver.PublishTask(&task); err != nil {
		log.Errorf("push task error: [%s]", err)
		return err
	}
	if task.Period > 0 {
		if err := d.smarter.PushTask(task.TaskId); err != nil {
			log.Errorf("push task to smarter error: [%s]", err)
			return err
		}
	}
	log.Info("ok")
	return nil
}

// push a job to broker and saver
// first save job in saver
// then push it to broker
func (d *dispatcher) PushJob(job core.Job) error {
	log := Logger.WithField("action", "PushJob").WithField("job", job)
	if err := d.saver.PublishJob(&job); err != nil {
		log.Errorf("push job to saver error: [%s]", err)
		return err
	}
	if err := d.broker.PushJob(&job); err != nil {
		log.Errorf("push job to broker error: [%s]", err)
		return err
	}
	log.Info("ok")
	return nil
}

// get a job and delete it from the queue
func (d *dispatcher) GetJob(queue string) (*core.Job, error) {
	var job *core.Job
	log := Logger.WithField("action", "GetJob").WithField("queue", queue)
	job, err := d.broker.GetJob(queue)
	if err != nil {
		log.Errorf("get job error: [%s]", err)
		return nil, err
	}
	log.WithField("job", job).Info("ok")
	return job, nil
}

package linda

import (
	"github.com/amlun/linda/linda/core"
	"github.com/twinj/uuid"
)

// push a [period] task to saver
// then push to smarter if it is a scheduled task
func (l *Linda) PushTask(task core.Task) error {
	log := Logger.WithField("action", "PushTask").WithField("task", task)
	if err := l.saver.SaveTask(&task); err != nil {
		log.Errorf("push task error: [%s]", err)
		return err
	}
	if task.Period > 0 {
		if err := l.smarter.PushTask(task.TaskId); err != nil {
			log.Errorf("push task to smarter error: [%s]", err)
			return err
		}
	}
	log.Info("ok")
	return nil
}

func (l *Linda) Schedule() (string, error) {
	Logger.Info("schedule task from smarter")
	return l.smarter.GetTask()
}

func (l *Linda) GetTask(taskId string) (*core.Task, error) {
	task, err := l.saver.GetTask(taskId)
	if err != nil {
		Logger.Errorf("get task from saver error: [%s]", err)
		return nil, err
	}
	return task, nil

}

func (l *Linda) ReSetTask(taskId string) error {
	if err := l.smarter.PushTask(taskId); err != nil {
		Logger.Errorf("push task to smarter error: [%s]", err)
		return err
	}
	return nil
}

// first save job in saver
// then push it to broker
func (l *Linda) PushJob(job core.Job) error {
	log := Logger.WithField("action", "PushJob").WithField("job", job)
	if err := l.saver.SaveJob(&job); err != nil {
		log.Errorf("push job to saver error: [%s]", err)
		return err
	}
	if err := l.broker.PushJob(&job); err != nil {
		log.Errorf("push job to broker error: [%s]", err)
		return err
	}
	log.Info("ok")
	return nil
}

// get a job and delete it from the queue
func (l *Linda) GetJob(queue string) (*core.Job, error) {
	var job *core.Job
	log := Logger.WithField("action", "GetJob").WithField("queue", queue)
	job, err := l.broker.GetJob(queue)
	if err != nil {
		log.Errorf("get job error: [%s]", err)
		return nil, err
	}
	log.WithField("job", job).Info("ok")
	return job, nil
}

func (l *Linda) ScheduleTask(task *core.Task) func() {
	return func() {
		var job core.Job
		job.JobId = uuid.NewV4().String()
		job.TaskId = task.TaskId
		job.Queue = task.Queue
		job.Payload = task.Payload
		l.PushJob(job)
	}
}

// get task list
func (l *Linda) TaskList(taskList *core.TaskList) error {
	return l.saver.TaskList(taskList)
}

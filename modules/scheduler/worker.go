package scheduler

import (
	cron "github.com/carlescere/scheduler"
	"sync"
)

type worker struct {
	process
	sync.RWMutex
	jobList map[string]*cron.Job
}

func newWorker() (*worker, error) {
	process, err := newProcess("worker")
	if err != nil {
		return nil, err
	}
	jm := make(map[string]*cron.Job)
	return &worker{
		process: *process,
		jobList: jm,
	}, nil
}

// TODO
// With multi workers, there maybe have some problems
func (w *worker) work(taskIds <-chan string) {
	defer func() {
		log.Debug("worker work stop")
	}()
	log.Debug("worker work start, receive tasks from chan")
	for taskId := range taskIds {
		log.WithField("task_id", taskId).Debug("receive a task success")
		task, err := Linda.GetTask(taskId)
		if err != nil {
			log.Error(err)
			return
		}
		job, ok := w.jobList[taskId]
		if ok {
			log.WithField("task_id", taskId).Debug("previous cron job quit")
			job.Quit <- true
		}
		if job, err = cron.Every(task.Period).Seconds().Run(Linda.ScheduleTask(task)); err == nil {
			log.WithField("task_id", taskId).Debug("new cron job start")
			w.jobList[taskId] = job
		}
	}
	return
}

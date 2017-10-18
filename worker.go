package linda

import (
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type worker struct {
	process
}

type workerFunc func(...interface{}) error

func newWorker(id string) (*worker, error) {
	process, err := newProcess("worker-" + id)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &worker{
		process: *process,
	}, nil
}

func (w *worker) work(jobIDs <-chan string, monitor *sync.WaitGroup) {
	monitor.Add(1)
	go func() {
		logrus.Debugf("worker {%s} start...", w)
		defer func() {
			logrus.Debugf("worker {%s} stop...", w)
			defer monitor.Done()
		}()
		for jobID := range jobIDs {
			job, err := saver.Get(jobID)
			if err != nil || job == nil {
				logrus.Errorf("saver.Get(%v) error {%s}", jobID, err)
				continue
			}
			if workerFunc, ok := workers[job.Payload.Class]; ok {
				w.run(job, workerFunc)
			} else {
				// TODO
				logrus.Errorf("no worker for job {%s}", job)
			}
			if err := w.nextRun(job); err != nil {
				logrus.Error(err)
			}
			saver.Put(job)
		}
	}()
}

func (w *worker) run(job *Job, workerFunc workerFunc) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("workerFunc(%s) error {%s}", job, r)
			job.State.Retries++
		}
	}()
	logrus.Debugf("run job {%s}", job)
	if err := workerFunc(job.Payload.Args...); err != nil {
		panic(err)
	}
	job.State.RunTimes++
	job.State.LastRunAt = time.Now()
	job.State.Retries = 0
}

func (w *worker) nextRun(job *Job) error {
	var err error
	if job.Period == 0 {
		if job.State.Retries >= job.Retry {
			err = broker.Delete(job.Queue, job.ID)
		} else {
			// retry after 300 seconds...
			err = broker.Release(job.Queue, job.ID, 300)
		}
	} else {
		err = broker.Release(job.Queue, job.ID, job.Period)
	}
	return err
}

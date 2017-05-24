package linda

import (
	"github.com/sirupsen/logrus"
	"time"
)

type poller struct {
	process
}

func newPoller() (*poller, error) {
	process, err := newProcess("poller")
	if err != nil {
		return nil, err
	}

	return &poller{
		process: *process,
	}, nil
}

func (p *poller) poll(queue string, life int64, interval time.Duration, quit <-chan bool) <-chan *Job {
	logrus.Debug("poller poll start...")
	jobs := make(chan *Job)
	go func() {
		for {
			brokerConn.MigrateExpiredJobs(queue)
			logrus.Debugf("sleep for migrate expire jobs %v...", interval)
			timeout := time.After(interval)
			select {
			case <-timeout:
			}
		}
	}()
	go func() {
		defer func() {
			close(jobs)
			logrus.Debug("poller poll stop...")
		}()
		for {
			select {
			case <-quit:
				return
			default:
				job, err := brokerConn.Pop(queue, life)
				if err != nil {
					logrus.Error(err)
					return
				}
				if job != nil {
					select {
					case jobs <- job:
					case <-quit:
						return
					}
				} else {
					logrus.Debugf("sleep for get job %v...", interval)
					timeout := time.After(interval)
					select {
					case <-quit:
						return
					case <-timeout:
					}
				}
			}
		}
	}()
	return jobs
}

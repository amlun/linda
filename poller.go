package linda

import (
	"time"
	"github.com/sirupsen/logrus"
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

func (p *poller) poll(queue string, timeout int64, interval time.Duration) <-chan string {
	jobIDs := make(chan string)
	go func() {
		logrus.Debugf("poller {%s} start...", p)
		defer func() {
			logrus.Debugf("poller {%s} stop...", p)
			close(jobIDs)
		}()
		for {
			select {
			case <-quit:
				return
			default:
				broker.MigrateExpiredJobs(queue)
				jobID, _ := broker.Reserve(queue, timeout)
				if jobID != "" {
					select {
					case jobIDs <- jobID:
					case <-quit:
						return
					}
				} else { // sleep ...
					sleep := time.After(interval)
					logrus.Debugf("poller sleep {%s}", interval)
					select {
					case <-quit:
						return
					case <-sleep:
					}
				}
			}
		}
	}()
	return jobIDs
}

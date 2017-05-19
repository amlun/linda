package scheduler

import (
	"errors"
	"sync"
	"time"
)

type poller struct {
	process
	sync.RWMutex
	tasksMap map[string]time.Time
}

func newPoller() (*poller, error) {
	process, err := newProcess("poller")
	if err != nil {
		return nil, err
	}
	// control the queue length
	m := make(map[string]time.Time)
	return &poller{
		process:  *process,
		tasksMap: m,
	}, nil
}

// fetch all taskIds
func (p *poller) poll(interval time.Duration, quit <-chan bool) <-chan string {
	taskIds := make(chan string)
	go func() {
		log.Debug("poller poll start")
		// re add task ids to smarter
		defer func() {
			defer close(taskIds)
			log.Debug("poller poll stop")
			p.flush()
		}()
		for {
			select {
			case <-quit:
				return
			default:
				taskId, err := p.getTask()
				if err == nil {
					select {
					case taskIds <- taskId:
					case <-quit:
						return
					}
				} else {
					log.WithField("interval", interval).Debug("sleeping for new task")
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
	return taskIds
}

// get task from smarter and enqueue task id to poller Map
func (p *poller) getTask() (string, error) {
	defer p.Unlock()
	p.Lock()
	// task num limit
	// limit()
	// get task from smarter
	taskId, err := Linda.Schedule()
	if err != nil {
		return "", err
	}
	if taskId == "" {
		return "", errors.New("task is empty")
	}
	log.WithField("task_id", taskId).Debug("get task from smarter")
	// set poller task map
	p.tasksMap[taskId] = time.Now()
	return taskId, nil
}

// flush all task id, re add to smarter
func (p *poller) flush() error {
	defer p.Unlock()
	p.Lock()
	log.Debug("flush poller's all tasks return to smarter")
	for taskId := range p.tasksMap {
		Linda.ReSetTask(taskId)
	}
	return nil
}

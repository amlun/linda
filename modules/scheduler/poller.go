package scheduler

import (
	"errors"
	"fmt"
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
		// re add task ids to smarter
		defer func() {
			defer close(taskIds)
			p.flush()
		}()
		for {
			select {
			case <-quit:
				return
			default:
				// get task from smarter
				fmt.Println("get task from smarter")
				taskId, err := Linda.Schedule()
				fmt.Println(taskId, err)
				if err != nil {
					fmt.Errorf("schedule error [%s]", err)
					return
				}
				err = p.enqueue(taskId)
				if err == nil {
					select {
					case taskIds <- taskId:
					case <-quit:
						return
					}
				} else {
					fmt.Printf("Sleeping for %v\n", interval)
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

// enqueue the task id to poller queue
func (p *poller) enqueue(taskId string) error {
	defer p.Unlock()
	p.Lock()
	if taskId == "" {
		return errors.New("task is empty")
	}
	_, ok := p.tasksMap[taskId]
	if ok {
		return errors.New("task already in poller")
	}
	p.tasksMap[taskId] = time.Now()
	return nil
}

// flush all task id, re add to smarter
func (p *poller) flush() error {
	defer p.Unlock()
	p.Lock()
	for taskId := range p.tasksMap {
		Linda.ReSetTask(taskId)
	}
	return nil
}

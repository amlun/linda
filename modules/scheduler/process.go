package scheduler

import (
	"fmt"
	"os"
)

type process struct {
	Hostname string
	Pid      int
	ID       string
}

func newProcess(id string) (*process, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &process{
		Hostname: hostname,
		Pid:      os.Getpid(),
		ID:       id,
	}, nil
}

func (p *process) String() string {
	return fmt.Sprintf("%s:%d-%s", p.Hostname, p.Pid, p.ID)
}

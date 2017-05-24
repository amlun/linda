package linda

import "fmt"

type Job struct {
	Queue   string `json:"queue"`
	Payload Payload
}

func (j *Job) String() string {
	return fmt.Sprintf("In queue: %s | handle: %s(%v)", j.Queue, j.Payload.Class, j.Payload.Args)
}

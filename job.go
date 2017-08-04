package linda

import "fmt"

// Job is the basic unit of this package
// it contains queue name and payload
type Job struct {
	Queue   string `json:"queue"`
	Period  int64 `json:"period"`
	Payload Payload
}

// String format job to string
func (j *Job) String() string {
	return fmt.Sprintf("In queue: %s | handle: %s(%v)", j.Queue, j.Payload.Class, j.Payload.Args)
}

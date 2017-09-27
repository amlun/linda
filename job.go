package linda

import (
	"time"
	"encoding/json"
)

// Job is the basic unit of this package
// it contains queue name and payload
type Job struct {
	ID        string    `json:"id"`
	Queue     string    `json:"queue"`
	Period    int64     `json:"period"`
	Retry     int64     `json:"retry"`
	CreatedAt time.Time `json:"created_at"`
	Payload   Payload   `json:"payload"`
	State     State     `json:"state"`
}

// Payload is the job's payload
type Payload struct {
	Class string        `json:"class"`
	Args  []interface{} `json:"args"`
}

// State is the job's running state
type State struct {
	RunTime time.Time `json:"run_time"`
	Retries int64     `json:"retries"`
}

// String format job to string
func (j *Job) String() string {
	bytes, err := json.Marshal(j)
	if err != nil {
		return ""
	}
	return string(bytes)
}

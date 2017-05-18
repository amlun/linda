package core

import "time"

type Job struct {
	JobId   string    `json:"job_id"`
	RunTime time.Time `json:"run_time"`
	Delay   int       `json:"delay"`
	Status  int       `json:"status"`
	TaskId  string    `json:"task_id"`
	Queue   string    `json:"queue"`
	Payload
}

//const (
//	PENDING = iota
//	STARTED
//	RETRY
//	SUCCESS
//	FAILURE
//)

package core

import "time"

type Job struct {
	JobId   string    `json:"job_id"`
	RunTime time.Time `json:"run_time"`
	Delay   int       `json:"delay"`
	Status  int       `json:"status"`
	Task
}

//const (
//	PENDING = iota
//	STARTED
//	RETRY
//	SUCCESS
//	FAILURE
//)

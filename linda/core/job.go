package core

import "time"

type Job struct {
	ID      string `json:"job_id"`
	RunTime time.Time `json:"run_time"`
	Delay   time.Duration `json:"delay"`
	TaskId  string `json:"task_id"`
	Func    string `json:"func"`
	Args    []string `json:"args"`
	Status  int `json:"status"`
}

//const (
//	PENDING = iota
//	STARTED
//	RETRY
//	SUCCESS
//	FAILURE
//)

package core

import "time"

type Task struct {
	ID        string `json:"task_id"`
	Frequency time.Duration `json:"frequency"`
	Func      string `json:"func"`
	Args      []string `json:"args"`
}

type TaskList struct {
	Total int `json:"total"`
	Tasks []Task `json:"tasks"`
	State string `json:"state"`
}

package core

type Task struct {
	TaskId string `json:"task_id"`
	Period int    `json:"period"`
	Queue  string `json:"queue"`
	Payload
}

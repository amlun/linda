package core

type Task struct {
	TaskId string   `json:"task_id"`
	Period int      `json:"period"`
	Func   string   `json:"func"`
	Args   []string `json:"args"`
}

type TaskList struct {
	Total int    `json:"total"`
	Tasks []Task `json:"tasks"`
	State string `json:"state"`
}

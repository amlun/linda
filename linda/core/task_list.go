package core

type TaskList struct {
	Total int    `json:"total"`
	Tasks []Task `json:"tasks"`
	State string `json:"state"`
}

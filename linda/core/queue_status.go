package core

type QueueStatus struct {
	Queue  string `json:"queue"`
	Length int    `json:"length"`
}

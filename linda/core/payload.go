package core

type Payload struct {
	Args []string `json:"args"`
	Func string   `json:"func"`
}

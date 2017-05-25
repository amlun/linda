package linda

// Payload is the job's body
type Payload struct {
	Class string        `json:"class"`
	Args  []interface{} `json:"args"`
}

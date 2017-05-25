package linda

var (
	workers map[string]workerFunc
)

func init() {
	workers = make(map[string]workerFunc)
}

// register worker with workerFunc
// map to the Job Payload.Class
func RegisterWorkers(class string, worker workerFunc) {
	workers[class] = worker
}

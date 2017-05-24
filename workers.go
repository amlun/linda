package linda

var (
	workers map[string]workerFunc
)

func init() {
	workers = make(map[string]workerFunc)
}

func RegisterWorkers(class string, worker workerFunc) {
	workers[class] = worker
}

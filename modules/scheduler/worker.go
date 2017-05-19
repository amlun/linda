package scheduler

//import (
//	"fmt"
//	cron "github.com/carlescere/scheduler"
//	"sync"
//)
//
//type worker struct {
//	process
//}
//
//func newWorker() (*worker, error) {
//	process, err := newProcess("worker")
//	if err != nil {
//		return nil, err
//	}
//	return &worker{
//		process: *process,
//	}, nil
//}
//
//func (w *worker) work(taskIds <-chan string, monitor *sync.WaitGroup) {
//	monitor.Add(1)
//	go func() {
//		defer monitor.Done()
//		for taskId := range taskIds {
//			fmt.Sprintf("schedule task [%s]", taskId)
//			task, err := Linda.GetTask(taskId)
//			if err != nil {
//				fmt.Errorf("work task error [%s]", err)
//				return
//			}
//			cron.Every(task.Period).Seconds().Run(Linda.ScheduleTask(task))
//		}
//	}()
//}

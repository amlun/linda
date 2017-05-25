// Package linda is a background manager to poll jobs from broker and dispatch them to multi workers.
//
// Linda Broker provides a unified API across different broker (queue) services.
//
// Brokers allow you to defer the processing of a time consuming task.
//
// Use ReleaseWithDelay func, you can implement a cron job service.
// Simple Useage:
// package main
//
// import (
//      "fmt"
//      "github.com/amlun/linda"
// )
//
// func init() {
//	settings := linda.Settings{
//		Queue:         "scheduler",
//		Connection:      "redis://localhost:6379/",
//		Timeout:       60,
//		IntervalFloat: 5.0,
//		Concurrency:   1,
//	}
//	linda.SetSettings(settings)
//	linda.RegisterWorkers("DispatcherSeed", DispatcherSeed)
// }
//
// func main() {
//	if err := linda.Run(); err != nil {
//		fmt.Println("Error:", err)
//	}
// }
//
// func DispatcherSeed(job *linda.Job) error {
//	broker := linda.GetBroker()
//	// get seed info
//	// do seed job
//	// release job with delay (like a cron job)
//	broker.ReleaseWithDelay("scheduler", job, 60)
//	return nil
// }

package linda

package main

import (
	"fmt"
	"github.com/amlun/linda"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	linda.RegisterWorkers("printArgs", PrintArgs)
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	// broker
	b, _ := linda.NewBroker("redis://localhost:6379/")
	// saver
	s, _ := linda.NewSaver("redis://localhost:6379/")
	// config
	c := linda.Config{
		Queue:     "test",
		Timeout:   60,
		Interval:  time.Second,
		WorkerNum: 4,
	}
	quit := signals()
	linda.Init(c, b, s)
	go func() {
		defer func() {
			linda.Quit()
		}()
		<-quit
	}()

	if err := linda.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}

func PrintArgs(args ...interface{}) error {
	fmt.Println(args)
	return nil
}

// Signal Handling
func signals() <-chan bool {
	quit := make(chan bool)
	go func() {
		signals := make(chan os.Signal)
		defer close(signals)
		signal.Notify(signals, syscall.SIGQUIT, syscall.SIGTERM, os.Interrupt)
		defer signalStop(signals)
		<-signals
		quit <- true
	}()
	return quit
}

// Stops signals channel.
func signalStop(c chan<- os.Signal) {
	signal.Stop(c)
}

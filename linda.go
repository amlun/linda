package linda

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var (
	brokerConn  Broker
	initMutex   sync.Mutex
	initialized bool
	settings    Settings
)

// SetSettings set linda settings
func SetSettings(sets Settings) {
	settings = sets
}

// Init linda with settings
// get instance of broker
func Init() error {
	initMutex.Lock()
	defer initMutex.Unlock()
	if !initialized {
		logrus.Debug("start linda and init...")
		if err := flags(); err != nil {
			return err
		}
		logrus.Infof("init the linda with redis connection: %s", settings.Connection)
		logrus.Infof("init the linda with job queue name: %s", settings.Queue)
		logrus.Infof("init the linda with job timeout: %d seconds", settings.Timeout)
		logrus.Infof("init the linda with worker interval: %.1f", settings.IntervalFloat)
		logrus.Infof("init the linda with worker numbers: %d", settings.Concurrency)
		b, err := NewBroker(settings.Connection)
		if err != nil {
			return err
		}
		brokerConn = b
		initialized = true
	}
	return nil
}

// Close linda with close broker
func Close() {
	initMutex.Lock()
	defer initMutex.Unlock()
	if initialized {
		logrus.Infof("close linda...")
		brokerConn.Close()
		initialized = false
	}
}

// GetBroker will return the initialized broker of linda
func GetBroker() Broker {
	if initialized {
		return brokerConn
	}
	return nil
}

// Run linda, it also call init function self
func Run() error {
	err := Init()
	if err != nil {
		return err
	}
	defer Close()
	logrus.Debug("start the poller to get jobs and migrate expired jobs")
	poller, err := newPoller()
	if err != nil {
		return err
	}
	quit := signals()
	jobs := poller.poll(settings.Queue, settings.Timeout, time.Duration(settings.Interval), quit)
	var monitor sync.WaitGroup
	logrus.Debugf("start %d workers to do the job", settings.Concurrency)
	for i := 0; i < settings.Concurrency; i++ {
		worker, err := newWorker(strconv.Itoa(i))
		if err != nil {
			return err
		}
		worker.work(jobs, &monitor)
	}
	monitor.Wait()
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

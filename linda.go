package linda

import (
	"github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"errors"
)

var (
	broker            Broker
	saver             Saver
	initMutex         sync.Mutex
	initialized       bool
	config            *Config
	quit              chan bool
	ErrNotInitialized = errors.New("you must init linda first")
)

// Open linda with config
// get instance of broker
func Init(c Config, b Broker, s Saver) error {
	initMutex.Lock()
	defer initMutex.Unlock()
	if !initialized {
		logrus.Debugf("init linda...")
		config = &c
		quit = make(chan bool)
		// init the broker
		broker = b
		// init the saver
		saver = s
		// set initialized true
		initialized = true
	}
	return nil
}

// Close linda with close broker and saver
func Close() {
	initMutex.Lock()
	defer initMutex.Unlock()
	if initialized {
		logrus.Debugf("close linda...")
		broker.Close()
		saver.Close()
		initialized = false
	}
}

func Quit() {
	close(quit)
}

// Run linda, it also call init function self
func Run() error {
	if !initialized {
		return ErrNotInitialized
	}
	defer Close()
	if err := run(); err != nil {
		return err
	}
	return nil
}

func run() error {
	// poller
	poller, err := newPoller()
	if err != nil {
		return err
	}
	jobIDs := poller.poll(config.Queue, config.Timeout, config.Interval)

	// workers
	var monitor sync.WaitGroup
	for i := 0; i < config.WorkerNum; i++ {
		worker, err := newWorker(strconv.Itoa(i))
		if err != nil {
			return err
		}
		worker.work(jobIDs, &monitor)
	}
	monitor.Wait()

	return nil
}

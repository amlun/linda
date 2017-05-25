// Running linda
//
// Flags
//
// -queue=default
// -- worker run jobs of the queue's name
//
// -interval=5.0
// -- polling job wait time if there is no job in the queue
//
// -concurrency=4
// --executing workers number, every worker running in it's own goroutine
//
// -timeout=60
// --life of the job in reserved state
// if the job is out of the time, poller will kick it back off the ready queue
//
// -connection=redis://localhost:6379/
// --the url of the broker
// redis://user:pass@host:port/db or unix:///path/to/redis.sock

package linda

import (
	"errors"
	"flag"
	"github.com/sirupsen/logrus"
)

func init() {
	flag.StringVar(&settings.Queue, "queue", "default", "queue name")
	flag.Float64Var(&settings.IntervalFloat, "interval", 1.0, "sleep interval(seconds) when no jobs are found")
	flag.IntVar(&settings.Concurrency, "concurrency", 2, "the maximum number of concurrently workers")
	flag.Int64Var(&settings.Timeout, "timeout", 0, "the reserved job life time(seconds), 0 means job will be deleted directly")
	flag.StringVar(&settings.Connection, "connection", "redis://localhost:6379/", "the url of the broker connection")
}

func flags() error {
	if !flag.Parsed() {
		logrus.Debug("parse the flag")
		flag.Parse()
	}

	if err := settings.Interval.SetFloat(settings.IntervalFloat); err != nil {
		return err
	}

	if settings.IntervalFloat < 1 {
		return errors.New("interval must gte 1, poll too quickly")
	}

	return nil
}

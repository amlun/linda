package linda

import (
	"flag"
	"github.com/sirupsen/logrus"
)

func init() {
	flag.StringVar(&settings.Queue, "queues", "", "queue name")
	flag.Float64Var(&settings.IntervalFloat, "interval", 5.0, "sleep interval when no jobs are found")
	flag.IntVar(&settings.Concurrency, "concurrency", 2, "the maximum number of concurrently workers")
	flag.StringVar(&settings.Connection, "connection", "redis://localhost:6379/", "the url of the broker connection")
	flag.BoolVar(&settings.Ack, "ack", false, "if the job should be reserved for result ack")
}

func flags() error {
	if !flag.Parsed() {
		logrus.Debug("parse the flag")
		flag.Parse()
	}

	if err := settings.Interval.SetFloat(settings.IntervalFloat); err != nil {
		return err
	}

	return nil
}

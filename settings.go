package linda

import (
	"fmt"
	"strconv"
	"time"
)

// linda's settings
type Settings struct {
	Queue         string
	Connection    string
	IntervalFloat float64
	Timeout       int64
	Interval      intervalFlag
	Concurrency   int
}

type intervalFlag time.Duration

// set interval flag with string value
func (i *intervalFlag) Set(value string) error {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	i.SetFloat(f)
	return nil
}

// set interval flag with float value
func (i *intervalFlag) SetFloat(value float64) error {
	*i = intervalFlag(time.Duration(value * float64(time.Second)))
	return nil
}

// interval flag to string
func (i *intervalFlag) String() string {
	return fmt.Sprint(*i)
}

package redis

import (
	"fmt"
	"github.com/amlun/linda/linda/smarter"
	"testing"
)

func TestGetTask(t *testing.T) {
	s, err := smarter.NewSmarter("redis://10.60.81.83:6379")
	if err != nil {
		fmt.Println(err)
	}
	task, err := s.GetTask()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(task)
}

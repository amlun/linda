package cassandra

import (
	"fmt"
	"github.com/amlun/linda/linda/saver"
	"testing"
)

func TestGetTask(t *testing.T) {
	s, err := saver.NewSaver("cassandra://cassandra:cassandra@10.60.81.83:9042/linda")
	if err != nil {
		fmt.Println(err)
	}

	task, err := s.GetTask("b02c7bc8-a3c9-42ff-9c2a-8c7fca22d0ad")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(task)
}

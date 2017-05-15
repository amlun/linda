package cassandra

import (
	"encoding/base64"
	"fmt"
	neturl "net/url"
	"testing"
)

func TestScheduleTask(t *testing.T) {
	var saver Saver
	url, _ := neturl.Parse("cassandra://cassandra:cassandra@10.60.81.83:9042/linda")
	saver.Connect(url)

	err := saver.ScheduleTask("abc")
	if err != nil {
		fmt.Errorf("publishe timing task err [%s]", err)
	}
}

func TestTaskList(t *testing.T) {
	//var saver Saver
	//url, _ := neturl.Parse("cassandra://cassandra:cassandra@10.60.81.83:9042/linda")
	//saver.Connect(url)
	//
	//taskList, _ := saver.TaskList("")
	//fmt.Println(taskList)
	//
	//for taskList.State != "" {
	//	taskList, _ := saver.TaskList(taskList.State)
	//	fmt.Println(taskList)
	//}

	fmt.Println(base64.URLEncoding.DecodeString(""))
}

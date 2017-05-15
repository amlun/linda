package main

import (
	"github.com/amlun/linda/linda"
	"github.com/amlun/linda/modules/monitor"
)

func main() {
	var config = linda.Config{
		BrokerURL: "redis://10.60.81.83:6379",
		SaverURL:  "cassandra://cassandra:cassandra@10.60.81.83:9042/linda",
	}
	l := linda.NewLinda(&config)
	m := monitor.New(l)
	m.Start()
}

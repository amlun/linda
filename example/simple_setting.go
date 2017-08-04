package main

import (
	"fmt"
	"github.com/amlun/linda"
)

var queue = "simple"

// RPUSH simple "{\"queue\":\"simple\",\"period\":0,\"Payload\":{\"class\":\"MyClass\",\"args\":[\"a\",\"b\",\"c\"]}}"
// RPUSH simple "{\"queue\":\"simple\",\"period\":0,\"Payload\":{\"class\":\"MyClass\",\"args\":[\"x\",\"y\",\"z\"]}}"
// RPUSH simple "{\"queue\":\"simple\",\"period\":0,\"Payload\":{\"class\":\"MyClass\",\"args\":[1,2,3]}}"
func init() {
	settings := linda.Settings{
		Queue:         queue,
		Connection:    "redis://localhost:6379/",
		IntervalFloat: 5.0,
		Timeout:       0,
		Concurrency:   4,
	}
	linda.SetSettings(settings)
	linda.RegisterWorkers("MyClass", MyFunc)
}

func main() {
	if err := linda.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}

func MyFunc(args ...interface{}) error {
	fmt.Println(args)
	return nil
}

package main

import (
	"fmt"
	"github.com/amlun/linda"
)

func init() {
	linda.RegisterWorkers("PrintArgs", PrintArgs)
}

func main() {
	if err := linda.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}

func PrintArgs(job *linda.Job) error {
	fmt.Println(job.Payload.Args)
	return nil
}

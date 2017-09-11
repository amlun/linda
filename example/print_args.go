package main

import (
	"fmt"
	"github.com/amlun/linda"
)

func init() {
	linda.RegisterWorkers("printArgs", PrintArgs)
}

func main() {
	if err := linda.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}

func PrintArgs(args ...interface{}) error {
	fmt.Println(args)
	return nil
}

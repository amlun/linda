package main

import (
	"fmt"
	"github.com/amlun/linda"
)

// RPUSH print "{\"queue\":\"print\",\"Payload\":{\"class\":\"printArgs\",\"args\":[\"a\",\"b\",\"c\"]}}"
// RPUSH print "{\"queue\":\"print\",\"Payload\":{\"class\":\"printArgs\",\"args\":[\"A\",\"B\",\"C\"]}}"
// RPUSH print "{\"queue\":\"print\",\"Payload\":{\"class\":\"printArgs\",\"args\":[1,2,3,4,5,6,7]}}"
func init() {
	linda.RegisterWorkers("PrintArgs", PrintArgs)
}

// go run print_args -queue=print -connection=redis://localhost:6379/
func main() {
	if err := linda.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}

func PrintArgs(args ...interface{}) error {
	fmt.Println(args)
	return nil
}

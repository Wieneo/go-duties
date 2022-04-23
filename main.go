package main

import (
	"fmt"

	"github.com/Wieneo/go-duties/v2/pkg/duties"
)

func main() {
	manager := duties.NewDutyManager()

	task, err := manager.TaskList.AddTask("test", nil, test)
	if err != nil {
		fmt.Println(err)
	}

	depTask, err := manager.TaskList.GetTask("test")
	if err != nil {
		fmt.Println(err)
	}

	if err := task.AddDependency(depTask); err != nil {
		fmt.Println(err)
	}
}

var debug = 1

func test(data interface{}) error {
	fmt.Printf("test called\n")
	fmt.Printf("Debug: %d\n", debug)
	return nil
}

package main

import (
	"errors"
	"fmt"

	"github.com/Wieneo/go-duties/v2/pkg/duties"
)

func main() {
	manager := duties.NewDutyManager()

	task, err := manager.TaskList.AddTask("test", test)
	if err != nil {
		fmt.Println(err)
	}

	task2, err := manager.TaskList.AddTask("test2", test)
	if err != nil {
		fmt.Println(err)
	}

	task.AddDependency(task2)

	manager.Execute()

	fmt.Println(task.GetStatus().State)
	fmt.Println(task2.GetStatus().State)
}

var debug = 1

func test(data interface{}) error {
	fmt.Printf("test called\n")
	fmt.Printf("Debug: %d\n", debug)
	return errors.New("oof")
}

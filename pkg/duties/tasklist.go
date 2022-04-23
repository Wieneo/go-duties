package duties

import (
	"strings"

	"github.com/Wieneo/go-duties/v2/pkg/duties/utils"
)

type TaskList struct {
	tasks []Task
}

func (tl *TaskList) GetTask(name string) (*Task, error) {
	for i, k := range tl.tasks {
		if strings.EqualFold(k.name, name) {
			return &tl.tasks[i], nil
		}
	}
	return nil, TaskNotFound
}

func (tl *TaskList) AddTask(name string, dependencies []*Task, call func(data interface{}) error) (*Task, error) {
	if utils.IsEmpty(name) {
		return nil, EmptyTaskName
	}

	if call == nil {
		return nil, NoCallFunction
	}

	for _, k := range tl.tasks {
		if strings.EqualFold(k.name, name) {
			return nil, DuplicateTask
		}
	}

	initDependencies := make([]*Task, 0)
	if dependencies != nil {
		initDependencies = dependencies
	}

	tl.tasks = append(tl.tasks, Task{
		name:         name,
		dependencies: initDependencies,
		call:         call,
		taskList:     tl,
	})
	return tl.GetTask(name)
}

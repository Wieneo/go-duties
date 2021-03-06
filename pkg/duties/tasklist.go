package duties

import (
	"strings"

	"github.com/Wieneo/go-duties/v2/pkg/duties/utils"
)

type TaskList struct {
	Tasks []*Task
}

func (tl *TaskList) GetTask(name string) (*Task, error) {
	for i, k := range tl.Tasks {
		if strings.EqualFold(k.name, name) {
			return tl.Tasks[i], nil
		}
	}
	return nil, ErrTaskNotFound
}

func (tl *TaskList) AddTask(name string, call func(data interface{}) error, preflight func(data interface{}) error, data interface{}) (*Task, error) {
	if utils.IsEmpty(name) {
		return nil, ErrEmptyTaskName
	}

	if call == nil {
		return nil, ErrNoCallFunction
	}

	for _, k := range tl.Tasks {
		if strings.EqualFold(k.name, name) {
			return nil, ErrDuplicateTask
		}
	}

	initDependencies := make([]*Task, 0)

	newTask := Task{
		name:         name,
		dependencies: initDependencies,
		call:         call,
		taskList:     tl,
		preflight:    preflight,
		status: TaskStatus{
			State: TaskStateCreated,
		},
		data: data,
	}

	tl.Tasks = append(tl.Tasks, &newTask)
	return tl.GetTask(name)
}

package duties

import (
	"strings"
	"time"
)

type Task struct {
	name         string
	taskList     *TaskList
	dependencies []*Task
	call         func(data interface{}) error
	preflight    func(data interface{}) error
	status       TaskStatus
}

type TaskStatus struct {
	PreflightStartTime time.Time
	PrefLightEndTime   time.Time
	StartTime          time.Time
	EndTime            time.Time
	Error              error
	State              TaskState
}

type TaskState int

const (
	TaskStateCreated TaskState = iota
	TaskStatePending
	TaskStateInPreflight
	TaskStatePreflightSucceeded
	TaskStateRunning
	TaskStateSucceded
	TaskStateFailed
	TaskStateDependencyFailed
	TaskStatePreFlightFailed
)

func (t *Task) GetName() string {
	return t.name
}

func (t *Task) GetStatus() TaskStatus {
	return t.status
}

func (t *Task) AddDependency(dep *Task) error {
	//Can't have ourselfs as dependency
	if strings.EqualFold(t.name, dep.name) {
		return DependencySelfReference
	}

	//Check if its already an dependency
	for _, k := range t.dependencies {
		if strings.EqualFold(k.name, dep.name) {
			return DuplicateDependency
		}
	}

	//Check if tasks exists in our tasklist
	found := false
	for _, k := range t.taskList.tasks {
		if strings.EqualFold(k.name, dep.name) {
			found = true
		}
	}

	if !found {
		return TaskNotFound
	}

	t.dependencies = append(t.dependencies, dep)
	return nil
}

func (t *Task) runPreflight(data interface{}) {
	t.status.PreflightStartTime = time.Now()
	t.status.State = TaskStateInPreflight

	if err := t.preflight(data); err != nil {
		t.status.State = TaskStatePreFlightFailed
		t.status.Error = err
	} else {
		t.status.State = TaskStatePreflightSucceeded
	}
	t.status.PrefLightEndTime = time.Now()
}

func (t *Task) execute(data interface{}) {
	t.status.StartTime = time.Now()
	t.status.State = TaskStateRunning

	if err := t.call(data); err != nil {
		t.status.State = TaskStateFailed
		t.status.Error = err
	} else {
		t.status.State = TaskStateSucceded
	}
	t.status.EndTime = time.Now()
}

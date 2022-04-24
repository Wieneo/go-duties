package duties

import (
	"strings"
	"sync"
	"time"
)

type Task struct {
	mu           sync.Mutex
	name         string
	taskList     *TaskList
	dependencies []*Task
	call         func(data interface{}) error
	preflight    func(data interface{}) error
	data         interface{}
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

type TaskState string

const (
	TaskStateCreated            TaskState = "Created"
	TaskStatePending            TaskState = "Pending"
	TaskStateInPreflight        TaskState = "InPreflight"
	TaskStatePreflightSucceeded TaskState = "PreflightSucceeded"
	TaskStateRunning            TaskState = "Running"
	TaskStateSucceded           TaskState = "Succeeded"
	TaskStateFailed             TaskState = "Failed"
	TaskStateDependencyFailed   TaskState = "DependencyFailed"
	TaskStatePreFlightFailed    TaskState = "PreflightFailed"
)

func (t *Task) GetName() string {
	return t.name
}

func (t *Task) GetStatus() TaskStatus {
	return t.status
}

func (t *Task) setStatus(status TaskState) {
	if t.status.State != status {
		t.status.State = status
		logInfo(`Task "%s" is now in state: %s`, t.name, status)
	}
}

func (t *Task) AddDependency(dep *Task) error {
	//Can't have ourselfs as dependency
	if strings.EqualFold(t.name, dep.name) {
		return ErrDependencySelfReference
	}

	//Check if its already an dependency
	for _, k := range t.dependencies {
		if strings.EqualFold(k.name, dep.name) {
			return ErrDuplicateDependency
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
		return ErrTaskNotFound
	}

	t.dependencies = append(t.dependencies, dep)
	return nil
}

func (t *Task) runPreflight(data interface{}) {
	//Lock task object as this function will be a seperate go-routine most of the time
	t.mu.Lock()
	defer t.mu.Unlock()

	t.status.PreflightStartTime = time.Now()
	t.setStatus(TaskStateInPreflight)

	if t.preflight != nil {
		if err := t.preflight(data); err != nil {
			t.setStatus(TaskStatePreFlightFailed)
			t.status.Error = err
		} else {
			t.setStatus(TaskStatePreflightSucceeded)
		}
	} else {
		t.setStatus(TaskStatePreflightSucceeded)
	}
	t.status.PrefLightEndTime = time.Now()
}

func (t *Task) runCall(data interface{}) {
	//Lock task object as this function will be a seperate go-routine most of the time
	t.mu.Lock()
	defer t.mu.Unlock()

	t.status.StartTime = time.Now()
	t.setStatus(TaskStateRunning)

	if err := t.call(data); err != nil {
		t.setStatus(TaskStateFailed)
		t.status.Error = err
	} else {
		t.setStatus(TaskStateSucceded)
	}
	t.status.EndTime = time.Now()
}

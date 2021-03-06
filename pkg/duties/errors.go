package duties

import "errors"

var ErrDuplicateDependency error = errors.New("duplicate dependency")
var ErrDependencyChangeForbidden error = errors.New("dependencies can't be modified after execution was started")
var ErrDuplicateTask error = errors.New("task name needs to be unique")
var ErrEmptyTaskName error = errors.New("task name may not be empty")
var ErrNoCallFunction error = errors.New("task must implement a call function")
var ErrTaskNotFound = errors.New("task not found")
var ErrDependencySelfReference = errors.New("task can't have itself as depdendency")
var ErrDependencyLoop error = errors.New("task dependencies prevent any task from running")
var ErrDependencyFailed error = errors.New("one or more of the tasks dependencies failed")

package duties

import "errors"

var DuplicateDependency error = errors.New("duplicate dependency")
var DuplicateTask error = errors.New("task name needs to be unique")
var EmptyTaskName error = errors.New("task name may not be empty")
var NoCallFunction error = errors.New("task must implement a call function")
var TaskNotFound = errors.New("task not found")
var DependencySelfReference = errors.New("task can't have itself as depdendency")

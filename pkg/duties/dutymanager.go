package duties

import "fmt"

type DutyManager struct {
	TaskList TaskList
	TaskData interface{}
}

func NewDutyManager() DutyManager {
	return DutyManager{
		TaskList: TaskList{
			tasks: make([]*Task, 0),
		},
	}
}

func (dm *DutyManager) Execute() {
	tl := &dm.TaskList

	//Set all tasks to pending
	for i := range tl.tasks {
		tl.tasks[i].status.State = TaskStatePending
		fmt.Println(tl.tasks[i].dependencies)
	}

	completed := false
	for !completed {

		taskToBeDone := tasksInState(tl, TaskStatePending)
		if len(taskToBeDone) > 0 {
			for i := range taskToBeDone {
				task := taskToBeDone[i]

				allDependenciesCompleted := true
				for _, k := range task.dependencies {
					if k.status.State != TaskStateSucceded {
						allDependenciesCompleted = false
					}

					if k.status.State == TaskStateFailed {
						task.status.State = TaskStateDependencyFailed
					}
				}

				if allDependenciesCompleted {
					task.execute(dm.TaskData)
				}
			}
		} else {
			completed = true
		}
	}
}

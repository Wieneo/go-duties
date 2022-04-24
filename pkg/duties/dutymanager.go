package duties

import (
	"time"
)

type DutyManager struct {
	TaskList TaskList
}

func NewDutyManager() DutyManager {
	return DutyManager{
		TaskList: TaskList{
			tasks: make([]*Task, 0),
		},
	}
}

func (dm *DutyManager) Execute() error {
	//Disable logging for dry-run
	loggingDisabled = true

	err := dm.runTasks(true)

	//Re-enable logging after dry run and before maybe returning
	loggingDisabled = false

	if err != nil {
		return err
	}

	return dm.runTasks(false)
}

func (dm *DutyManager) runTasks(dryRun bool) error {
	tl := &dm.TaskList

	//Set all tasks to pending
	for i := range tl.tasks {
		tl.tasks[i].setStatus(TaskStatePending)
	}

	completed := false
	for !completed {
		loopStart := time.Now()

		taskToBeDone := tasksWaiting(tl)
		if len(taskToBeDone) > 0 {
			tasksRunInThisWave := 0
			tasksDoingSomething := len(tasksDoingSomething(tl))

			for i := range taskToBeDone {
				task := taskToBeDone[i]

				allDependenciesCompleted := true
				for _, k := range task.dependencies {
					if k.status.State != TaskStateSucceded {
						allDependenciesCompleted = false
					}

					if k.status.State == TaskStateFailed || k.status.State == TaskStatePreFlightFailed {
						task.setStatus(TaskStateDependencyFailed)
					}
				}

				if allDependenciesCompleted {
					if !dryRun {
						if task.GetStatus().State == TaskStatePending {
							task.setStatus(TaskStateInPreflight)
							go task.runPreflight(task.data)
						}

						if task.GetStatus().State == TaskStatePreflightSucceeded {
							task.setStatus(TaskStateRunning)
							go task.runCall(task.data)
						}
					} else {
						task.setStatus(TaskStateSucceded)
					}

					tasksRunInThisWave++
				}
			}

			//Check if we are stuck (ring dependencies or something like that)
			if tasksRunInThisWave == 0 && tasksDoingSomething == 0 {
				logError("Executing failed because of dependency loop", ErrDependencyLoop)
				return ErrDependencyLoop
			}

		} else {
			//Check if there as still tasks running
			if len(tasksDoingSomething(tl)) == 0 {
				completed = true
			}
		}

		//Throttle loop execution
		executionDuration := time.Since(loopStart)
		if executionDuration < 1*time.Millisecond {
			time.Sleep(1*time.Millisecond - executionDuration)
		}
	}

	return nil
}

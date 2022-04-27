package duties

import (
	"time"
)

type DutyManager struct {
	TaskList TaskList
}

//NewDutyManager returns a new duty manager instance
func NewDutyManager() DutyManager {
	return DutyManager{
		TaskList: TaskList{
			Tasks: make([]*Task, 0),
		},
	}
}

func (dm *DutyManager) PendingTasks() []*Task {
	return tasksInState(&dm.TaskList, TaskStatePending)
}

func (dm *DutyManager) RunningTasks() []*Task {
	return tasksDoingSomething(&dm.TaskList)
}

func (dm *DutyManager) SuccededTasks() []*Task {
	return tasksInState(&dm.TaskList, TaskStateSucceded)
}

func (dm *DutyManager) FailedTasks() []*Task {
	tasks := append(tasksInState(&dm.TaskList, TaskStatePreFlightFailed), tasksInState(&dm.TaskList, TaskStateFailed)...)
	return append(tasks, tasksInState(&dm.TaskList, TaskStateDependencyFailed)...)
}

//Execute starts the execution of all defined tasks
//This method also performs a dry run, in order to check if all dependencies are valid and no dependency loop was created
//This method will block the current thread until all tasks were processed
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
	for i := range tl.Tasks {
		tl.Tasks[i].setStatus(TaskStatePending)
	}

	completed := false
	for !completed {
		loopStart := time.Now()

		//Get all tasks in state Pending / PreflightSucceded
		//These states indicate we need to start preflight or execution
		taskToBeDone := tasksWaiting(tl)
		if len(taskToBeDone) > 0 {
			//tasksRunInThisWave gets incremented for every task that enters preflight or execution during this iteration
			tasksRunInThisWave := 0

			for i := range taskToBeDone {
				task := taskToBeDone[i]

				//allDependenciesCompleted gets set to false if a dependency hasn't finished execution
				//We will retry execution of this tasks next iteration
				allDependenciesCompleted := true
				for _, k := range task.dependencies {
					if k.status.State != TaskStateSucceded {
						allDependenciesCompleted = false
					}

					//If one of the tasks dependencies failed, we can't execute this task
					if k.status.State == TaskStateFailed || k.status.State == TaskStatePreFlightFailed {
						task.setStatus(TaskStateDependencyFailed)
						task.status.Error = ErrDependencyFailed
						break
					}
				}

				if allDependenciesCompleted {
					//Don't do anything on a dry run -> Just set the task's state to succeded
					if !dryRun {
						if task.GetStatus().State == TaskStatePending {
							//Set status here because setting the state via the spawned go routine is too slow and the next iteration will see the old status
							task.setStatus(TaskStateInPreflight)
							go task.runPreflight(task.data)
						}

						if task.GetStatus().State == TaskStatePreflightSucceeded {
							//Set status here because setting the state via the spawned go routine is too slow and the next iteration will see the old status
							task.setStatus(TaskStateRunning)
							go task.runCall(task.data)
						}
					} else {
						task.setStatus(TaskStateSucceded)
					}

					tasksRunInThisWave++

				}
			}

		} else {
			//Check if there as still tasks running
			//Block current thread until we successfully processed all tasks
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

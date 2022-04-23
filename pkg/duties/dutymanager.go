package duties

type DutyManager struct {
	TaskList TaskList
}

func NewDutyManager() DutyManager {
	return DutyManager{
		TaskList: TaskList{
			tasks: make([]Task, 0),
		},
	}
}

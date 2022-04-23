package duties

func tasksInState(tl *TaskList, state TaskState) []*Task {
	tasks := make([]*Task, 0)
	for i, k := range tl.tasks {
		if k.status.State == state {
			tasks = append(tasks, tl.tasks[i])
		}
	}
	return tasks
}

package tomeit

import "time"

type task struct {
	id                     int64
	user                   *user
	title                  string
	expectedPomodoroNumber int
	dueOn                  time.Time
	isCompleted            bool
	completedAt            time.Time
	createdAt              time.Time
	updatedAt              time.Time
}

func hasUserTask(db dbInterface, taskID int64, user *user) bool {
	task, err := db.getTaskByID(taskID)
	if err != nil {
		return false
	}

	if task.user.id == user.id {
		return true
	} else {
		return false
	}
}

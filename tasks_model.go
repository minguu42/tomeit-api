package tomeit

import "time"

type task struct {
	id        int64
	user      *user
	name      string
	priority  int
	deadline  time.Time
	isDone    bool
	createdAt time.Time
	updatedAt time.Time
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

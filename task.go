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

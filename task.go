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

type Task struct {
	ID                  int
	UserID              int
	User                User
	Title               string
	ExpectedPomodoroNum int
	DueAt               *time.Time
	IsCompleted         bool
	CompletedAt         *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

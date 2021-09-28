package tomeit

import "time"

type Task struct {
	ID                  int
	UserID              int
	User                User
	Title               string
	ExpectedPomodoroNum int
	DueOn               *time.Time
	IsCompleted         bool
	CompletedOn         *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

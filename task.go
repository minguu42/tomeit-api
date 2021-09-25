package tomeit

import "time"

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

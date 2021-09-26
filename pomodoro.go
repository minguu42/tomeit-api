package tomeit

import "time"

type Pomodoro struct {
	ID        int
	UserID    int
	User      User
	TaskID    int
	Task      Task
	CreatedAt time.Time
}

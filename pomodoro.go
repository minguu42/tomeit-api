package tomeit

import "time"

type pomodoro struct {
	id          int64
	user        *user
	task        *task
	completedAt time.Time
	createdAt   time.Time
}

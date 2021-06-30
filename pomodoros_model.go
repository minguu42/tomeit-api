package tomeit

import "time"

type pomodoroLog struct {
	id        int64
	user      *user
	task      *task
	createdAt time.Time
}

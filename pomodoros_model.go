package tomeit

import "time"

type pomodoroRecord struct {
	id        int64
	user      *user
	task      *task
	createdAt time.Time
}

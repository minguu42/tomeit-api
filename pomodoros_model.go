package tomeit

import "time"

type pomodoroLog struct {
	id        int64
	userId    int64
	task      *Task
	createdAt time.Time
}

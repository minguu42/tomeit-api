package tomeit

import "time"

type Task struct {
	id        int64
	userId    int64
	name      string
	priority  int
	deadline  time.Time
	isDone    bool
	createdAt time.Time
	updatedAt time.Time
}

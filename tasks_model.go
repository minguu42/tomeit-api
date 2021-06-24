package tomeit

import "time"

type Task struct {
	id        int
	userId    int
	name      string
	priority  int
	deadline  time.Time
	isDone    bool
	createdAt time.Time
	updatedAt time.Time
}

package tomeit

import "time"

type task struct {
	id        int64
	user      *user
	name      string
	priority  int
	deadline  time.Time
	isDone    bool
	createdAt time.Time
	updatedAt time.Time
}

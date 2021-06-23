package tomeit

import "time"

type Task struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
	Deadline string `json:"deadline,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

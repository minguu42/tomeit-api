package tomeit

type Task struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
	Deadline string `json:"deadline"`
}

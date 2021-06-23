package tomeit

import (
	"fmt"
)

func createTask(userId int, name string, priority int, deadline string) (int64, error) {
	const q = `
INSERT INTO tasks (user_id, name, priority, deadline)
VALUES (?, ?, ?, ?);
`
	r, err := db.Exec(q, userId, name, priority, deadline)
	if err != nil {
		return 0, fmt.Errorf("exec failed: %w", err)
	}

	id, err := r.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("lastInsertId failed: %w", err)
	}

	return id, nil
}

func getTaskById(id int64) (Task, error) {
	const q = `SELECT * FROM tasks WHERE id = ?`

	var t Task

	if err := db.QueryRow(q, id).Scan(&t.ID, &t.UserID, &t.Name, &t.Priority, &t.Deadline, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return Task{}, fmt.Errorf("queryRow failed: %w", err)
	}

	return t, nil
}

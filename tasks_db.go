package tomeit

import (
	"fmt"
	"time"
)

func (db DB) createTask(userID int64, name string, priority int, deadline time.Time) (int64, error) {
	const q = `INSERT INTO tasks (user_id, name, priority, deadline) VALUES (?, ?, ?, ?)`

	r, err := db.Exec(q, userID, name, priority, deadline.Format("2006-01-02"))
	if err != nil {
		return 0, fmt.Errorf("exec failed: %w", err)
	}

	id, err := r.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("lastInsertId failed: %w", err)
	}

	return id, nil
}

func (db DB) getTaskByID(id int64) (*task, error) {
	const q = `SELECT name, priority, deadline, is_done, created_at, updated_at FROM tasks WHERE id = ?`

	t := task{id: id}

	if err := db.QueryRow(q, id).Scan(&t.name, &t.priority, &t.deadline, &t.isDone, &t.createdAt, &t.updatedAt); err != nil {
		return nil, fmt.Errorf("queryRow failed: %w", err)
	}

	return &t, nil
}

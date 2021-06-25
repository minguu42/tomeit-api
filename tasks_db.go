package tomeit

import (
	"fmt"
	"time"
)

func createTask(userId int64, name string, priority int, deadline time.Time) (int64, error) {
	const q = `
INSERT INTO tasks (user_id, name, priority, deadline)
VALUES (?, ?, ?, ?);
`
	r, err := db.Exec(q, userId, name, priority, deadline.Format("2006-01-02"))
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

	if err := db.QueryRow(q, id).Scan(&t.id, &t.userId, &t.name, &t.priority, &t.deadline, &t.isDone, &t.createdAt, &t.updatedAt); err != nil {
		return Task{}, fmt.Errorf("queryRow failed: %w", err)
	}

	return t, nil
}

func getUndoneTasksByUserID(userId int64) ([]*Task, error) {
	const q = `
SELECT * FROM tasks
WHERE user_id = ? AND is_done = FALSE
`
	var tasks []*Task
	rows, err := db.Query(q, userId)
	if err != err {
		return nil, fmt.Errorf("getUndoneTasks failed: %w", err)
	}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.id, &t.name, &t.priority, &t.deadline, &t.createdAt, &t.updatedAt); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		tasks = append(tasks, &t)
	}
	return tasks, nil
}

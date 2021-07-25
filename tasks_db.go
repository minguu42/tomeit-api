package tomeit

import (
	"fmt"
	"time"
)

func (db *DB) createTask(userID int64, name string, priority int, deadline time.Time) (int64, error) {
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

func (db *DB) getTaskByID(id int64) (*task, error) {
	const q = `
SELECT T.name, T.priority, T.deadline, T.is_done, T.created_at, T.updated_at, U.id, U.digest_uid 
FROM tasks AS T
JOIN users AS U ON T.user_id = U.id
WHERE T.id = ?
`

	var u user
	t := task{id: id, user: &u}

	if err := db.QueryRow(q, id).Scan(&t.name, &t.priority, &t.deadline, &t.isDone, &t.createdAt, &t.updatedAt, &u.id, &u.digestUID); err != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	return &t, nil
}

func (db *DB) getUndoneTasksByUser(user *user) ([]*task, error) {
	const q = `
SELECT id, name, priority, deadline, is_done, created_at, updated_at FROM tasks
WHERE user_id = ? AND is_done = FALSE
ORDER BY updated_at
LIMIT 30
`
	var ts []*task
	rows, err := db.Query(q, user.id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	for rows.Next() {
		t := task{
			user: user,
		}
		if err := rows.Scan(&t.id, &t.name, &t.priority, &t.deadline, &t.isDone, &t.createdAt, &t.updatedAt); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		ts = append(ts, &t)
	}

	return ts, nil
}

func (db *DB) getDoneTasksByUser(user *user) ([]*task, error) {
	const q = `
SELECT id, name, priority, deadline, is_done, created_at, updated_at FROM tasks
WHERE user_id = ? AND is_done = TRUE
ORDER BY updated_at
LIMIT 30
`
	var tasks []*task
	rows, err := db.Query(q, user.id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	for rows.Next() {
		t := task{
			user: user,
		}
		if err := rows.Scan(&t.id, &t.name, &t.priority, &t.deadline, &t.isDone, &t.createdAt, &t.updatedAt); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		tasks = append(tasks, &t)
	}

	return tasks, nil
}

func (db *DB) getPomodoroCountByID(id int64) (int, error) {
	const q = `SELECT COUNT(*) FROM pomodoro_logs WHERE task_id = ?`

	var c int
	if err := db.QueryRow(q, id).Scan(&c); err != nil {
		return 0, fmt.Errorf("scan failed: %w", err)
	}

	return c, nil
}

func (db *DB) doneTask(taskID int64) error {
	const q = `UPDATE tasks SET is_done = TRUE WHERE id = ?`

	_, err := db.Exec(q, taskID)
	if err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}

	return nil
}

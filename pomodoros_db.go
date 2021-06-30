package tomeit

import (
	"fmt"
)

func (db *DB) createPomodoroLog(userID, taskID int64) (int64, error) {
	const q = `INSERT INTO pomodoro_logs (user_id, task_id) VALUES (?, ?)`

	r, err := db.Exec(q, userID, taskID)
	if err != nil {
		return 0, fmt.Errorf("exec failed: %w", err)
	}

	id, err := r.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("lastInsertId failed: %w", err)
	}

	return id, nil
}

func (db *DB) getPomodoroLogByID(id int64) (*pomodoroLog, error) {
	const q = `
SELECT P.created_at, U.id, U.digest_uid, T.id, T.name, T.priority, T.deadline, T.is_done, T.created_at, T.updated_at
FROM pomodoro_logs AS P
JOIN users AS U ON P.user_id = U.id
JOIN tasks AS T ON P.task_id = T.id
WHERE P.id = ?
`

	var u user
	var t task
	p := pomodoroLog{
		id:   id,
		user: &u,
		task: &t,
	}
	if err := db.QueryRow(q, id).Scan(&p.createdAt, &u.id, &u.digestUID, &t.id, &t.name, &t.priority, &t.deadline, &t.isDone, &t.createdAt, &t.updatedAt); err != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	return &p, nil
}

func (db *DB) getPomodoroLogsByUser(user *user) ([]*pomodoroLog, error) {
	const q = `
SELECT P.id, P.created_at, T.id, T.name, T.priority, T.deadline, T.is_done, T.created_at, T.updated_at
FROM pomodoro_logs AS P
JOIN tasks AS T ON P.task_id = T.id
WHERE P.user_id = ?
ORDER BY P.created_at
LIMIT 30
`
	var ps []*pomodoroLog

	rows, err := db.Query(q, user.id)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	for rows.Next() {
		var t task
		p := pomodoroLog{
			user: user,
		}
		if err := rows.Scan(&p.id, &p.createdAt, &t.id, &t.name, &t.priority, &t.deadline, &t.isDone, &t.createdAt, &t.updatedAt); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		p.task = &t
		ps = append(ps, &p)
	}

	return ps, nil
}

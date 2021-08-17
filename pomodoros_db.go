package tomeit

import (
	"database/sql"
	"fmt"
	"time"
)

func (db *DB) createPomodoro(userID, taskID int64) (int64, error) {
	const q = `INSERT INTO pomodoros (user_id, task_id) VALUES (?, ?)`

	r, err := db.Exec(q, userID, taskID)
	if err != nil {
		return 0, fmt.Errorf("db.Exec failed: %w", err)
	}

	id, err := r.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("result.LastInsertId failed: %w", err)
	}

	return id, nil
}

func (db *DB) getPomodoroByID(id int64) (*pomodoro, error) {
	const q = `
SELECT P.completed_at, P.created_at, U.id, U.digest_uid, T.id, T.title, T.expected_pomodoro_number, T.due_on, T.is_completed, T.completed_at,T.created_at, T.updated_at
FROM pomodoros AS P
JOIN users AS U ON P.user_id = U.id
JOIN tasks AS T ON P.task_id = T.id
WHERE P.id = ?
`

	var u user
	var t task
	var nullDueOn sql.NullTime
	var nullCompletedAt sql.NullTime
	p := pomodoro{
		id:   id,
		user: &u,
		task: &t,
	}
	if err := db.QueryRow(q, id).Scan(&p.completedAt, &p.createdAt, &u.id, &u.digestUID, &t.id, &t.title, &t.expectedPomodoroNumber, &nullDueOn, &t.isCompleted, &nullCompletedAt, &t.createdAt, &t.updatedAt); err != nil {
		return nil, fmt.Errorf("row.Scan failed: %w", err)
	}
	t.dueOn = nullDueOn.Time
	t.completedAt = nullCompletedAt.Time

	return &p, nil
}

type getPomodorosOptions struct {
	existCompletedOn bool
	completedOn      time.Time
}

func (db *DB) getPomodorosByUser(user *user, options *getPomodorosOptions) ([]*pomodoro, error) {
	var optionList []string
	if options != nil {
		if options.existCompletedOn {
			optionList = append(optionList, "AND DATE(P.completed_at) = '"+options.completedOn.Format("2006-01-02")+"'")
		}
	}

	q := `
SELECT P.id, P.completed_at, P.created_at, T.id, T.title, T.expected_pomodoro_number, T.due_on, T.is_completed, T.completed_at,T.created_at, T.updated_at
FROM pomodoros AS P
JOIN tasks AS T ON P.task_id = T.id
WHERE P.user_id = ?
`
	for _, option := range optionList {
		q = q + option
	}
	q = q + `
ORDER BY P.created_at
LIMIT 30
`
	var ps []*pomodoro
	rows, err := db.Query(q, user.id)
	if err != nil {
		return nil, fmt.Errorf("db.Query failed: %w", err)
	}

	for rows.Next() {
		var t task
		var nullDueOn sql.NullTime
		var nullCompletedAt sql.NullTime
		p := pomodoro{
			user: user,
		}
		if err := rows.Scan(&p.id, &p.completedAt, &p.createdAt, &t.id, &t.title, &t.expectedPomodoroNumber, &nullDueOn, &t.isCompleted, &nullCompletedAt, &t.createdAt, &t.updatedAt); err != nil {
			return nil, fmt.Errorf("rows.Scan failed: %w", err)
		}
		t.dueOn = nullDueOn.Time
		t.completedAt = nullCompletedAt.Time
		p.task = &t
		ps = append(ps, &p)
	}

	return ps, nil
}

func (db *DB) getTodayPomodoroCount(user *user) (int, error) {
	today := time.Now().UTC().Format("2006-01-02")

	const q = `
SELECT COUNT(*) FROM pomodoro_logs
WHERE user_id = ? AND DATE(created_at) = ?
`
	var c int

	if err := db.QueryRow(q, user.id, today).Scan(&c); err != nil {
		return 0, fmt.Errorf("row.Scan failed: %w", err)
	}
	return c, nil
}

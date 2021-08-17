package tomeit

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"
)

func (db *DB) createTask(userID int64, name string, expectedPomodoroNumber int, dueOn time.Time) (int64, error) {
	const q = `INSERT INTO tasks (user_id, title, expected_pomodoro_number, due_on) VALUES (?, ?, ?, ?)`

	r, err := db.Exec(q, userID, name, expectedPomodoroNumber, dueOn.Format("2006-01-02"))
	if err != nil {
		return 0, fmt.Errorf("db.Exec failed: %w", err)
	}

	id, err := r.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("result.lastInsertId failed: %w", err)
	}

	return id, nil
}

func (db *DB) getTaskByID(id int64) (*task, error) {
	const q = `
SELECT T.title, T.expected_pomodoro_number, T.due_on, T.is_completed, T.completed_at, T.created_at, T.updated_at, U.id
FROM tasks AS T
JOIN users AS U ON T.user_id = U.id
WHERE T.id = ?
`

	var u user
	var nullDueOn sql.NullTime
	var nullCompletedAt sql.NullTime
	t := task{id: id, user: &u}

	if err := db.QueryRow(q, id).Scan(&t.title, &t.expectedPomodoroNumber, &nullDueOn, &t.isCompleted, &nullCompletedAt, &t.createdAt, &t.updatedAt, &u.id); err != nil {
		return nil, fmt.Errorf("row.Scan failed: %w", err)
	}

	t.dueOn = nullDueOn.Time
	t.completedAt = nullCompletedAt.Time

	return &t, nil
}

type getTasksOptions struct {
	existIsCompleted bool
	isCompleted      bool
	existCompletedOn bool
	completedOn      time.Time
}

func (db *DB) getTasksByUser(user *user, options *getTasksOptions) ([]*task, error) {
	var optionList []string
	if options != nil {
		if options.existIsCompleted {
			optionList = append(optionList, "AND is_completed = "+strconv.FormatBool(options.isCompleted))
		}
		if options.existCompletedOn {
			optionList = append(optionList, "AND DATE(completed_at) = '"+options.completedOn.Format("2006-01-02")+"'")
		}
	}

	q := `
SELECT id, title, expected_pomodoro_number, due_on, is_completed, completed_at, created_at, updated_at
FROM tasks
WHERE user_id = ?
`
	for _, option := range optionList {
		q = q + option
	}
	q = q + `
ORDER BY created_at
LIMIT 30
`
	var ts []*task
	rows, err := db.Query(q, user.id)
	if err != nil {
		return nil, fmt.Errorf("db.Query failed: %w", err)
	}

	for rows.Next() {
		t := task{
			user: user,
		}

		var nullDueOn sql.NullTime
		var nullCompletedAt sql.NullTime
		if err := rows.Scan(&t.id, &t.title, &t.expectedPomodoroNumber, &nullDueOn, &t.isCompleted, &nullCompletedAt, &t.createdAt, &t.updatedAt); err != nil {
			return nil, fmt.Errorf("rows.Scan failed: %w", err)
		}
		t.dueOn = nullDueOn.Time
		t.completedAt = nullCompletedAt.Time

		ts = append(ts, &t)
	}

	return ts, nil
}

func (db *DB) getActualPomodoroNumberByID(id int64) (int, error) {
	const q = `SELECT COUNT(*) FROM pomodoros WHERE task_id = ?`

	var c int
	if err := db.QueryRow(q, id).Scan(&c); err != nil {
		return 0, fmt.Errorf("row.Scan failed: %w", err)
	}

	return c, nil
}

type updateTaskOptions struct {
	existIsCompleted bool
}

func (db *DB) updateTask(task *task, options *updateTaskOptions) error {
	if options == nil {
		return errors.New("options must not be nil")
	}

	var optionList []string
	if options.existIsCompleted {
		optionList = append(optionList, "is_completed = "+strconv.FormatBool(task.isCompleted))
		now := time.Now()
		optionList = append(optionList, "completed_at = '"+now.Format("2006-01-02 15:04:05")+"'")
	}

	q := `UPDATE tasks SET`
	for i, option := range optionList {
		if i == 0 {
			q = q + " " + option + " "
		} else {
			q = q + ", " + option + " "
		}
	}
	q = q + `WHERE id = ?`

	_, err := db.Exec(q, task.id)
	if err != nil {
		return fmt.Errorf("db.Exec failed: %w", err)
	}
	return nil
}

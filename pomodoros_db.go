package tomeit

import (
	"fmt"
)

func createPomodoroLog(userId int64, taskId int64) (int64, error) {
	const q = `INSERT INTO pomodoro_logs (user_id, task_id) VALUES (?, ?)`

	r, err := db.Exec(q, userId, taskId)
	if err != nil {
		return 0, fmt.Errorf("exec failed: %w", err)
	}

	id, err := r.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("lastInsertId failed: %w", err)
	}

	return id, nil
}

func getPomodoroLogById(taskId int64) (*pomodoroLog, error) {
	const q = `
SELECT P.id, P.user_id, T.id, T.user_id, T.name, T.priority, T.deadline, T.is_done, T.created_at, T.updated_at, P.created_at
FROM pomodoro_logs AS P
JOIN tasks AS T ON P.task_id = T.id
WHERE P.task_id = ?
`
	var t Task
	var p pomodoroLog

	if err := db.QueryRow(q, taskId).Scan(&p.id, &p.userId, &t.id, &t.userId, &t.name, &t.priority, t.deadline, t.isDone, t.createdAt, t.updatedAt, p.createdAt); err != nil {
		return nil, fmt.Errorf("queryRow failed: %w", err)
	}

	p.task = &t
	return &p, nil
}

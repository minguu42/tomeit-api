package tomeit

import (
	"fmt"
	"time"
)

func (db *DB) createPomodoro(userID, taskID int) (int, error) {
	pomodoro := Pomodoro{
		UserID: userID,
		TaskID: taskID,
	}
	if err := db.Create(&pomodoro).Error; err != nil {
		return 0, fmt.Errorf("db.Create failed: %w", err)
	}
	return pomodoro.ID, nil
}

func (db *DB) getPomodoroByID(id int) (*Pomodoro, error) {
	var pomodoro Pomodoro

	if err := db.First(&pomodoro, id).Error; err != nil {
		return nil, fmt.Errorf("db.First failed: %w", err)
	}

	var t Task
	if err := db.First(&t, pomodoro.TaskID).Error; err != nil {
		return nil, fmt.Errorf("db.First failed: %w", err)
	}
	pomodoro.Task = t

	return &pomodoro, nil
}

type getPomodorosOptions struct {
	existCompletedOn bool
	completedOn      time.Time
}

func (db *DB) getPomodorosByUser(user *User, options *getPomodorosOptions) ([]Pomodoro, error) {
	const q = `
SELECT P.id, P.created_at, T.id, T.title, T.expected_pomodoro_num, T.due_at, T.is_completed, T.completed_at, T.created_at, T.updated_at
FROM pomodoros AS P
JOIN tasks AS T ON P.task_id = T.id
WHERE P.user_id = ?
ORDER BY P.created_at
LIMIT 30
`

	//if options != nil {
	//	if options.existCompletedOn {
	//		y, m, d := options.completedOn.Date()
	//		start := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	//		end := time.Date(y, m, d, 23, 59, 59, 0, time.UTC)
	//		q = q.Where("completed_at BETWEEN ? AND ?", start, end)
	//	}
	//}

	rows, err := db.Raw(q, user.ID).Rows()
	if err != nil {
		return nil, fmt.Errorf("db.Rows failed: %w", err)
	}

	var pomodoros []Pomodoro
	for rows.Next() {
		var p Pomodoro
		if err := rows.Scan(&p.ID, &p.CreatedAt, &p.Task.ID, &p.Task.Title, &p.Task.ExpectedPomodoroNum, &p.Task.DueAt, &p.Task.IsCompleted, &p.Task.CompletedAt, &p.Task.CreatedAt, &p.Task.UpdatedAt); err != nil {
			return nil, fmt.Errorf("rows.Scan failed: %w", err)
		}
		pomodoros = append(pomodoros, p)
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("rows.Close faield")
	}

	return pomodoros, nil
}

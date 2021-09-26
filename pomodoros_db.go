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
	q := db.Select("pomodoros.id, pomodoros.user_id, pomodoros.task_id, tasks.id, tasks.title, tasks.expected_pomodoro_num, tasks.due_at, tasks.is_completed, tasks.completed_at, tasks.created_at, tasks.updated_at, pomodoros.created_at").Joins("JOIN tasks ON pomodoros.task_id = tasks.id").Where("user_id = ?", user.ID)

	if options != nil {
		if options.existCompletedOn {
			y, m, d := options.completedOn.Date()
			start := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
			end := time.Date(y, m, d, 23, 59, 59, 0, time.UTC)
			q = q.Where("completed_at BETWEEN ? AND ?", start, end)
		}
	}

	var pomodoros []Pomodoro
	if err := q.Order("created_at").Limit(30).Find(&pomodoros).Error; err != nil {
		return nil, fmt.Errorf("db.Find failed: %w", err)
	}

	return pomodoros, nil
}

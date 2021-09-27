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

	if err := db.Joins("Task").First(&pomodoro, id).Error; err != nil {
		return nil, fmt.Errorf("db.First failed: %w", err)
	}

	return &pomodoro, nil
}

type getPomodorosOptions struct {
	completedOnExists bool
	completedOn       time.Time
}

func (db *DB) getPomodorosByUser(user *User, options *getPomodorosOptions) ([]Pomodoro, error) {
	q := db.Joins("Task").Where("pomodoros.user_id = ?", user.ID).Order("pomodoros.created_at").Limit(30)

	if options != nil {
		if options.completedOnExists {
			y, m, d := options.completedOn.Date()
			start := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
			end := time.Date(y, m, d, 23, 59, 59, 0, time.UTC)
			q = q.Where("pomodoros.completed_at BETWEEN ? AND ?", start, end)
		}
	}

	var pomodoros []Pomodoro
	if err := q.Find(&pomodoros).Error; err != nil {
		return nil, fmt.Errorf("q.Find failed: %w", err)
	}
	return pomodoros, nil
}

func (db *DB) deletePomodoro(pomodoro *Pomodoro) {
	db.Delete(pomodoro)
}

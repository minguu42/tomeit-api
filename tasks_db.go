package tomeit

import (
	"fmt"
	"time"
)

type taskDBInterface interface {
	createTask(userID int, title string, priority int, dueAt time.Time) (*Task, error)
	getTaskByID(id int) (*Task, error)
	getTasksByUser(user *User, options *getTasksOptions) ([]Task, error)
	getActualPomodoroNumByID(id int) (int, error)
	updateTask(task *Task) error
	deleteTask(task *Task) error
}

func (db *DB) createTask(userID int, title string, expectedPomodoroNum int, dueAt time.Time) (*Task, error) {
	now := time.Now()
	task := Task{
		UserID:              userID,
		Title:               title,
		ExpectedPomodoroNum: expectedPomodoroNum,
		DueAt:               &dueAt,
		IsCompleted:         false,
		CompletedAt:         nil,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	q := db.DB
	if dueAt.IsZero() {
		q = q.Omit("DueAt")
	}

	if err := q.Create(&task).Error; err != nil {
		return nil, fmt.Errorf("db.Create failed: %w", err)
	}
	return &task, nil
}

func (db *DB) getTaskByID(id int) (*Task, error) {
	var t Task

	if err := db.First(&t, id).Error; err != nil {
		return nil, fmt.Errorf("db.First failed: %w", err)
	}

	return &t, nil
}

type getTasksOptions struct {
	isCompletedExists bool
	isCompleted       bool
	completedOnExists bool
	completedOn       time.Time
}

func (db *DB) getTasksByUser(user *User, options *getTasksOptions) ([]Task, error) {
	q := db.Where("user_id = ?", user.ID)

	if options != nil {
		if options.isCompletedExists {
			q = q.Where("is_completed = ?", options.isCompleted)
		}
		if options.completedOnExists {
			y, m, d := options.completedOn.Date()
			start := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
			end := time.Date(y, m, d, 23, 59, 59, 0, time.UTC)
			q = q.Where("completed_at BETWEEN ? AND ?", start, end)
		}
	}

	var tasks []Task

	if err := q.Order("created_at").Limit(30).Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("db.Find failed: %w", err)
	}

	return tasks, nil
}

func (db *DB) getActualPomodoroNumByID(id int) (int, error) {
	var c int64

	if err := db.Model(&Pomodoro{}).Where("task_id = ?", id).Count(&c).Error; err != nil {
		return 0, fmt.Errorf("db.Count failed: %w", err)
	}

	return int(c), nil
}

func (db *DB) updateTask(task *Task) error {
	if err := db.Save(task).Error; err != nil {
		return fmt.Errorf("db.Save() failed: %w", err)
	}

	return nil
}

func (db *DB) deleteTask(task *Task) error {
	if err := db.Delete(task).Error; err != nil {
		return fmt.Errorf("db.Delete failed: %w", err)
	}

	return nil
}

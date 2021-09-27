package tomeit

import (
	"fmt"
	"time"
)

func (db *DB) createTask(userID int, title string, expectedPomodoroNum int, dueAt time.Time) (int, error) {
	task := Task{
		UserID:              userID,
		Title:               title,
		ExpectedPomodoroNum: expectedPomodoroNum,
		DueAt:               &dueAt,
	}

	q := db.DB
	if dueAt.IsZero() {
		q = q.Omit("DueAt")
	}

	if err := q.Create(&task).Error; err != nil {
		return 0, fmt.Errorf("db.Create failed: %w", err)
	}
	return task.ID, nil
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

//func (db *DB) getActualPomodoroNumberByID(id int64) (int, error) {
//	const q = `SELECT COUNT(*) FROM pomodoros WHERE task_id = ?`
//
//	var c int
//	if err := db.QueryRow(q, id).Scan(&c); err != nil {
//		return 0, fmt.Errorf("row.Scan failed: %w", err)
//	}
//
//	return c, nil
//}
//

func (db *DB) updateTask(task *Task) {
	db.Save(task)
}

func (db *DB) deleteTask(task *Task) {
	db.Delete(task)
}

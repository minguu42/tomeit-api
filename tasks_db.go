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

	q := db
	if dueAt.IsZero() {
		q.Omit("DueAt")
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
	existIsCompleted bool
	isCompleted      bool
	existCompletedOn bool
	completedOn      time.Time
}

func (db *DB) getTasksByUser(user *User, options *getTasksOptions) ([]Task, error) {
	q := db.Where("user_id = ?", user.ID)

	if options != nil {
		if options.existIsCompleted {
			q = q.Where("is_completed = ?", options.isCompleted)
		}
		if options.existCompletedOn {
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
//type updateTaskOptions struct {
//	isCompletedExists bool
//}
//
//func (db *DB) updateTask(task *task, options *updateTaskOptions) error {
//	if options == nil {
//		return errors.New("options must not be nil")
//	}
//
//	var optionList []string
//	if options.isCompletedExists {
//		optionList = append(optionList, "is_completed = "+strconv.FormatBool(task.isCompleted))
//		now := time.Now()
//		optionList = append(optionList, "completed_at = '"+now.Format("2006-01-02 15:04:05")+"'")
//	}
//
//	q := `UPDATE tasks SET`
//	for i, option := range optionList {
//		if i == 0 {
//			q = q + " " + option + " "
//		} else {
//			q = q + ", " + option + " "
//		}
//	}
//	q = q + `WHERE id = ?`
//
//	_, err := db.Exec(q, task.id)
//	if err != nil {
//		return fmt.Errorf("db.Exec failed: %w", err)
//	}
//	return nil
//}

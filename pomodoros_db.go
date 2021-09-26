package tomeit

import (
	"fmt"
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
//
//type getPomodorosOptions struct {
//	existCompletedOn bool
//	completedOn      time.Time
//}
//
//func (db *DB) getPomodorosByUser(user *user, options *getPomodorosOptions) ([]*pomodoro, error) {
//	var optionList []string
//	if options != nil {
//		if options.existCompletedOn {
//			optionList = append(optionList, "AND DATE(P.completed_at) = '"+options.completedOn.Format("2006-01-02")+"'")
//		}
//	}
//
//	q := `
//SELECT P.id, P.completed_at, P.created_at, T.id, T.title, T.expected_pomodoro_number, T.due_on, T.is_completed, T.completed_at,T.created_at, T.updated_at
//FROM pomodoros AS P
//JOIN tasks AS T ON P.task_id = T.id
//WHERE P.user_id = ?
//`
//	for _, option := range optionList {
//		q = q + option
//	}
//	q = q + `
//ORDER BY P.created_at
//LIMIT 30
//`
//	var ps []*pomodoro
//	rows, err := db.Query(q, user.id)
//	if err != nil {
//		return nil, fmt.Errorf("db.Query failed: %w", err)
//	}
//
//	for rows.Next() {
//		var t task
//		var nullDueOn sql.NullTime
//		var nullCompletedAt sql.NullTime
//		p := pomodoro{
//			user: user,
//		}
//		if err := rows.Scan(&p.id, &p.completedAt, &p.createdAt, &t.id, &t.title, &t.expectedPomodoroNumber, &nullDueOn, &t.isCompleted, &nullCompletedAt, &t.createdAt, &t.updatedAt); err != nil {
//			return nil, fmt.Errorf("rows.Scan failed: %w", err)
//		}
//		t.dueOn = nullDueOn.Time
//		t.completedAt = nullCompletedAt.Time
//		p.task = &t
//		ps = append(ps, &p)
//	}
//
//	return ps, nil
//}

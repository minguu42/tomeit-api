package tomeit

import (
	"database/sql"
	"log"
	"time"
)

type dbInterface interface {
	createUser(digestUID string) (*user, error)
	getUserByDigestUID(digestUID string) (*user, error)
	decrementRestCount(user *user) error

	createTask(userID int64, name string, priority int, deadline time.Time) (int64, error)
	getTaskByID(id int64) (*task, error)
	getTasksByUser(user *user, options *getTasksOptions) ([]*task, error)
	getActualPomodoroNumberByID(id int64) (int, error)
	updateTask(task *task, options *updateTaskOptions) error

	createPomodoro(userID, taskID int64) (int64, error)
	getPomodoroByID(id int64) (*pomodoro, error)
	getPomodorosByUser(user *user, options *getPomodorosOptions) ([]*pomodoro, error)
}

type DB struct {
	*sql.DB
}

func OpenDB(driver, dsn string) *DB {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Fatalln("Open db failed:", err)
	}

	isDBReady := false
	failureTimes := 0
	for !isDBReady {
		err := db.Ping()
		if err == nil {
			isDBReady = true
		} else {
			log.Println("Ping db failed. try again.")
			time.Sleep(time.Second * 15)
			failureTimes += 1
		}

		if failureTimes >= 2 {
			log.Fatalln("Ping db failed:", err)
		}
	}

	return &DB{db}
}

func CloseDB(db *DB) {
	if err := db.Close(); err != nil {
		log.Fatalln("Close db failed:", err)
	}
}

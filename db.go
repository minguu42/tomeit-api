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
	getTasksByUser(user *user) ([]*task, error)
	getDoneTasksByUser(user *user) ([]*task, error)
	doneTask(taskID int64) error

	createPomodoroLog(userID, taskID int64) (int64, error)
	getPomodoroLogByID(id int64) (*pomodoroLog, error)
	getPomodoroLogsByUser(user *user) ([]*pomodoroLog, error)
}

type DB struct {
	*sql.DB
}

func OpenDB(driverName, databaseUrl string) *DB {
	db, err := sql.Open(driverName, databaseUrl)
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
			time.Sleep(time.Second * 15)
			failureTimes += 1
		}

		if failureTimes >= 3 {
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

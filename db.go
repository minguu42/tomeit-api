package tomeit

import (
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type dbInterface interface {
	createUser(digestUID string) (*User, error)
	getUserByDigestUID(digestUID string) (*User, error)
	//decrementRestCount(user *user) error

	createTask(userID int, title string, priority int, dueAt time.Time) (int, error)
	getTaskByID(id int) (*Task, error)
	getTasksByUser(user *User, options *getTasksOptions) ([]Task, error)
	//getActualPomodoroNumberByID(id int) (int, error)
	//updateTask(task *task, options *updateTaskOptions) error

	//createPomodoro(userID, taskID int64) (int64, error)
	//getPomodoroByID(id int64) (*pomodoro, error)
	//getPomodorosByUser(user *user, options *getPomodorosOptions) ([]*pomodoro, error)
}

type DB struct {
	*gorm.DB
}

func OpenDB(dsn string) *DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableAutomaticPing: true,
	})
	if err != nil {
		log.Fatal("Open db failed:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("db.DB failed:", err)
	}

	isDBReady := false
	failureTimes := 0
	for !isDBReady {
		err := sqlDB.Ping()
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
	sqlDB, err := db.DB.DB()
	if err != nil {
		log.Fatal("db.DB failed:", err)
	}

	if err := sqlDB.Close(); err != nil {
		log.Fatal("sqlDB.Close failed:", err)
	}
}

package tomeit

import (
	"database/sql"
	"log"
	"time"
)

type dbInterface interface {
	createUser(digestUID string) (*user, error)
	getUserByDigestUID(digestUID string) (*user, error)
	createTask(userID int64, name string, priority int, deadline time.Time) (int64, error)
	getTaskByID(id int64) (*task, error)
	getTasksByUser(user *user) ([]*task, error)
}

type DB struct {
	*sql.DB
}

func OpenDB(driverName, databaseUrl string) *DB {
	db, err := sql.Open(driverName, databaseUrl)
	if err != nil {
		log.Fatalln("Open db failed:", err)
	}
	return &DB{db}
}

func CloseDB(db *DB) {
	if err := db.Close(); err != nil {
		log.Fatalln("Close db failed:", err)
	}
}

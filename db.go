package tomeit

import (
	"database/sql"
	"log"
)

type dbInterface interface {
	createUser(digestUID string) (*user, error)
	getUserByDigestUID(digestUID string) (*user, error)
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

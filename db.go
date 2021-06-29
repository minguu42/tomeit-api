package tomeit

import (
	"database/sql"
	"log"
)

type dbInterface interface{}

type DB struct {
	*sql.DB
}

func OpenDB(driverName, databaseUrl string) *sql.DB {
	db, err := sql.Open(driverName, databaseUrl)
	if err != nil {
		log.Fatalf("Open db failed: %v\n", err)
	}
	return db
}

func CloseDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Fatalf("Close db failed: %v\n", err)
	}
}

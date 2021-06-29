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
		log.Fatalln("Open db failed:", err)
	}
	return db
}

func CloseDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Fatalln("Close db failed:", err)
	}
}

package tomeit

import (
	"database/sql"
	"log"
	"os"
)

var db *sql.DB

func OpenDb() {
	var err error
	db, err = sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
}

func CloseDb() {
	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
}

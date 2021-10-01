package tomeit

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type dbInterface interface {
	userDBInterface
	taskDBInterface
	pomodoroDBInterface
}

type DB struct {
	*gorm.DB
}

func OpenDB(dsn string) *DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			IgnoreRecordNotFoundError: true,
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableAutomaticPing: true,
		Logger:               newLogger,
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

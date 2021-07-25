package tomeit

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	_ "github.com/go-sql-driver/mysql"
)

var (
	testClient *http.Client
	testUrl    string
	testDB     *DB
)

func TestMain(m *testing.M) {
	firebaseApp := &firebaseAppMock{}

	testDB = OpenDB("mysql", "test:password@tcp(localhost:13306)/db_test?parseTime=true")
	defer CloseDB(testDB)

	r := chi.NewRouter()

	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(UserCtx(testDB, firebaseApp))

	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", PostTask(testDB))

		r.Get("/undone", GetTasksUndone(testDB))

		r.Route("/done", func(r chi.Router) {
			r.Get("/", GetTasksDone(testDB))
			r.Put("/{taskID}", PutTaskDone(testDB))
		})
	})
	r.Route("/pomodoros", func(r chi.Router) {
		r.Route("/records", func(r chi.Router) {
			r.Post("/", PostPomodoroRecord(testDB))
			r.Get("/", GetPomodoroRecords(testDB))
			r.Get("/count/today", GetTodayPomodoroCount(testDB))
		})

		r.Get("/rest-count", GetRestCount)
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	testUrl = ts.URL
	testClient = &http.Client{}

	m.Run()
}

func setupTestDB() {
	file, err := os.ReadFile(filepath.Join(".", "build", "setup.sql"))
	if err != nil {
		log.Fatal("os.ReadFile failed:", err)
	}
	queries := strings.Split(string(file), ";")

	for _, query := range queries {
		if query == "" {
			break
		}

		_, err := testDB.Exec(query)
		if err != nil {
			log.Fatal("db.Exec failed:", err)
		}
	}

	const createTestUser = `INSERT INTO users (digest_uid) VALUES ('a2c4ba85c41f186283948b1a54efacea04cb2d3f54a88d5826a7e6a917b28c5a')`

	if _, err := testDB.Exec(createTestUser); err != nil {
		log.Fatalln("createTestUser failed:", err)
	}
}

func shutdownTestDB() {
	const dropPomodoroLogsTable = `DROP TABLE IF EXISTS pomodoro_logs`
	const dropTasksTable = `DROP TABLE IF EXISTS tasks`
	const dropUsersTable = `DROP TABLE IF EXISTS users`

	if _, err := testDB.Exec(dropPomodoroLogsTable); err != nil {
		log.Fatalln("shutdownTestDB failed:", err)
	}
	if _, err := testDB.Exec(dropTasksTable); err != nil {
		log.Fatalln("shutdownTestDB failed:", err)
	}
	if _, err := testDB.Exec(dropUsersTable); err != nil {
		log.Fatalln("shutdownTestDB failed:", err)
	}
}

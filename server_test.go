package tomeit

import (
	"log"
	"net/http"
	"net/http/httptest"
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
	const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
    id              INT      NOT NULL AUTO_INCREMENT PRIMARY KEY,
    digest_uid      CHAR(64) NOT NULL,
    rest_count      INT      DEFAULT 4 NOT NULL CHECK ( 1 <= rest_count AND rest_count <= 4 )
)
`
	const createTasksTable = `
CREATE TABLE IF NOT EXISTS tasks (
    id         INT          NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id    INT          NOT NULL,
    name       VARCHAR(120) NOT NULL,
    priority   INT          DEFAULT 0 NOT NULL CHECK ( 0 <= priority AND priority <= 3 ),
    deadline   DATE         DEFAULT ('0001-01-01') NOT NULL,
    is_done    BOOLEAN      DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP    DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
)
`
	const createPomodoroLogsTable = `
CREATE TABLE IF NOT EXISTS pomodoro_logs (
    id         INT       NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id    INT       NOT NULL,
    task_id    INT       NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (task_id) REFERENCES tasks(id)
)
`
	const createTestUser = `INSERT INTO users (digest_uid) VALUES ('a2c4ba85c41f186283948b1a54efacea04cb2d3f54a88d5826a7e6a917b28c5a')`

	if _, err := testDB.Exec(createUsersTable); err != nil {
		log.Fatalln("createUsersTable failed:", err)
	}
	if _, err := testDB.Exec(createTasksTable); err != nil {
		log.Fatalln("createTasksTable failed:", err)
	}
	if _, err := testDB.Exec(createPomodoroLogsTable); err != nil {
		log.Fatalln("createPomodoroLogsTable failed:", err)
	}
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

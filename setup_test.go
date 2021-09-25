package tomeit

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

var (
	testClient *http.Client
	testUrl    string
	testDB     *DB
)

func TestMain(m *testing.M) {
	firebaseApp := &firebaseAppMock{}

	testDB = OpenDB("test:password@tcp(localhost:13306)/db_test?charset=utf8mb4&parseTime=true")
	defer CloseDB(testDB)

	r := chi.NewRouter()

	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(UserCtx(testDB, firebaseApp))

	Route(r, testDB)

	ts := httptest.NewServer(r)
	defer ts.Close()

	testUrl = ts.URL
	testClient = &http.Client{}

	m.Run()
}

func setupTestDB(tb testing.TB) {
	file, err := os.ReadFile(filepath.Join(".", "build", "create_tables.sql"))
	if err != nil {
		tb.Fatal("os.ReadFile failed:", err)
	}
	queries := strings.Split(string(file), ";")

	for _, query := range queries {
		if query == "" {
			break
		}

		testDB.Exec(query)
	}

	const createTestUser = `INSERT INTO users (digest_uid) VALUES ('a2c4ba85c41f186283948b1a54efacea04cb2d3f54a88d5826a7e6a917b28c5a')`

	testDB.Exec(createTestUser)
}

func teardownTestDB() {
	const dropPomodorosTable = `DROP TABLE IF EXISTS pomodoros`
	const dropTasksTable = `DROP TABLE IF EXISTS tasks`
	const dropUsersTable = `DROP TABLE IF EXISTS users`

	testDB.Exec(dropPomodorosTable)
	testDB.Exec(dropTasksTable)
	testDB.Exec(dropUsersTable)
}

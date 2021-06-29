package tomeit

import (
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
)

func TestMain(m *testing.M) {
	firebaseApp := &firebaseAppMock{}

	db := OpenDB("mysql", "test:password@tcp(localhost:13306)/db_test?parseTime=true")
	defer CloseDB(db)

	r := chi.NewRouter()

	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(UserCtx(db, firebaseApp))

	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", PostTask(db))
		r.Get("/", GetTasks(db))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	testUrl = ts.URL
	testClient = &http.Client{}

	m.Run()
}

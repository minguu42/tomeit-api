package tomeit

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

var (
	client *http.Client
	url    string
)

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open("mysql", "test:password@tcp(localhost:13306)/db_test?parseTime=true")
	if err != nil {
		log.Fatal("Open db failed:", err)
	}
	db.SetConnMaxLifetime(time.Second)

	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(mockUserCtx)
	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", PostTask)
		//	r.Get("/", GetUndoneTasks)
		//
		//	r.Route("/done", func(r chi.Router) {
		//		r.Get("/", GetDoneTasks)
		//		r.Put("/{taskId}", PutTaskDone)
		//	})
		//})
		//r.Route("/pomodoros", func(r chi.Router) {
		//	r.Route("/logs", func(r chi.Router) {
		//		r.Post("/", PostPomodoroLog)
		//	})
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	url = ts.URL
	client = &http.Client{}

	m.Run()

	if err := db.Close(); err != nil {
		log.Fatal("Close db failed:", err)
	}
}

package tomeit

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	_ "github.com/go-sql-driver/mysql"
)

var (
	client *http.Client
	url    string
)

func TestMain(m *testing.M) {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))

	var err error
	db, err = sql.Open("mysql", "test:password@tcp(localhost:13306)/db_test?parseTime=true")
	if err != nil {
		log.Fatal("Open db failed:", err)
	}

	r.Use(mockUserCtx)
	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", PostTask)
		r.Get("/", GetTasks)
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

func TestPostTask(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		reqBody := strings.NewReader(`{"name": "タスク1", "priority": 2, "deadline": "2021-12-31"}`)
		req, err := http.NewRequest("POST", url+"/tasks", reqBody)
		if err != nil {
			t.Errorf("Create request failed: %v", err)
		}
		req.Header.Add("Authorization", "someIdToken")

		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("Do request failed: %v", err)
		}

		bytes, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			t.Errorf("Read response body failed: %v", err)
		}

		var body taskResponse
		if err := json.Unmarshal(bytes, &body); err != nil {
			t.Errorf("Unmarshal json failed: %v", err)
		}

		if resp.StatusCode != 201 {
			t.Errorf("Status code should be 201, but %v", resp.StatusCode)
		}
		if body.Id <= 0 {
			t.Errorf("Id should be positive number, but %v", body.Id)
		}
		if body.Name != "タスク1" {
			t.Errorf("Name should be タスク1, but %v", body.Name)
		}
		if body.Priority != 2 {
			t.Errorf("Priority should be 2, but %v", body.Priority)
		}
		if body.Deadline != "2021-12-31" {
			t.Errorf("Deadline should be 2021-12-31, but %v", body.Deadline)
		}
		if body.IsDone != false {
			t.Errorf("IsDone should be false, but %v", body.IsDone)
		}
		if body.PomodoroCount != 0 {
			t.Errorf("PomodoroCount should be 0, but %v", body.PomodoroCount)
		}
		if body.CreatedAt == "" {
			t.Errorf("CreatedAt is empty")
		}
		if body.UpdatedAt == "" {
			t.Errorf("UpdatedAt is empty")
		}
	})
}

func TestGetTasks(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		req, err := http.NewRequest("GET", url+"/tasks", nil)
		if err != nil {
			t.Errorf("Create request failed: %v", err)
		}
		req.Header.Add("Authorization", "someIdToken")

		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("Do request failed: %v", err)
		}

		bytes, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			t.Errorf("Read response body failed: %v", err)
		}

		var body tasksResponse
		if err := json.Unmarshal(bytes, &body); err != nil {
			t.Errorf("Unmarshal json failed: %v", err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("Status code should be 200, but %v", resp.StatusCode)
		}
	})
}

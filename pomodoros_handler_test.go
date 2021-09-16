package tomeit

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func setupTestPostPomodoro(tb testing.TB) {
	const createTask1 = `INSERT INTO tasks (user_id, title, expected_pomodoro_number, due_on, is_completed) VALUES (1, 'タスク1', 0, '2018-12-31', false)`

	if _, err := testDB.Exec(createTask1); err != nil {
		tb.Fatal("setupTestPostPomodoroLog failed:", err)
	}
}

func TestPostPomodoro(t *testing.T) {
	t.Run("ポモドーロを記録する", func(t *testing.T) {
		setupTestDB(t)
		setupTestPostPomodoro(t)

		reqBody := strings.NewReader(`{"taskID": 1 }`)
		req, err := http.NewRequest("POST", testUrl+"/pomodoros", reqBody)
		if err != nil {
			t.Error("Create request failed:", err)
		}

		resp, err := testClient.Do(req)
		if err != nil {
			t.Error("Do request failed:", err)
		}

		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Error("Read response failed:", err)
		}
		if err := resp.Body.Close(); err != nil {
			t.Error("Close response failed:", err)
		}

		var body pomodoroResponse
		if err := json.Unmarshal(bytes, &body); err != nil {
			t.Error("Unmarshal json failed:", err)
		}

		if resp.StatusCode != 201 {
			t.Error("Status code should be 201, but", resp.StatusCode)
		}

		if body.ID != 1 {
			t.Error("ID should be 1, but", body.ID)
		}
		task := body.Task
		if task.ID != 1 {
			t.Error("ID should be 1, but", task.ID)
		}
		if task.Title != "タスク1" {
			t.Error("Title should be タスク2, but", task.Title)
		}
		if task.ExpectedPomodoroNumber != 0 {
			t.Error("ExpectedPomodoroNumber should be 1, but", task.ExpectedPomodoroNumber)
		}
		if task.ActualPomodoroNumber != 1 {
			t.Error("ActualPomodoroNumber should be 1, but", task.ActualPomodoroNumber)
		}
		if task.DueOn != "2018-12-31T00:00:00Z" {
			t.Error("DueOn should be 2018-12-31T00:00:00Z, but", task.DueOn)
		}
		if task.IsCompleted != false {
			t.Error("IsCompleted should be false, but", task.IsCompleted)
		}
		if task.CompletedAt != "0001-01-01T00:00:00Z" {
			t.Error("CompletedAt should be 0001-01-01T00:00:00Z, but", task.CompletedAt)
		}
		if task.CreatedAt == "0001-01-01T00:00:00Z" {
			t.Error("CreatedAt should not be 0001-01-01T00:00:00Z")
		}
		if task.UpdatedAt == "0001-01-01T00:00:00Z" {
			t.Error("UpdatedAt should not be 0001-01-01T00:00:00Z")
		}
		if body.CompletedAt == "0001-01-01T00:00:00Z" {
			t.Error("CompletedAt should not be 0001-01-01T00:00:00Z")
		}
		if body.CreatedAt == "0001-01-01T00:00:00Z" {
			t.Error("CreatedAt does not exist")
		}

		shutdownTestDB(t)
	})
}

func setupTestGetPomodoros(tb testing.TB) {
	const createTask1 = `INSERT INTO tasks (user_id, title, expected_pomodoro_number, is_completed) VALUES (1, 'タスク1', 0, false)`
	const createPomodoro1 = `INSERT INTO pomodoros (user_id, task_id, completed_at) VALUES (1, 1, '2021-08-31 01:02:03')`
	const createPomodoro2 = `INSERT INTO pomodoros (user_id, task_id, completed_at) VALUES (1, 1, '2021-09-01 06:07:08')`

	if _, err := testDB.Exec(createTask1); err != nil {
		tb.Fatal("create Task1 failed:", err)
	}
	if _, err := testDB.Exec(createPomodoro1); err != nil {
		tb.Fatal("create Pomodoro1 failed:", err)
	}
	time.Sleep(time.Second * 1)
	if _, err := testDB.Exec(createPomodoro2); err != nil {
		tb.Fatal("create Pomodoro2 failed:", err)
	}
}

func TestGetPomodoros(t *testing.T) {
	t.Run("ポモドーロ記録を一覧取得する", func(t *testing.T) {
		setupTestDB(t)
		setupTestGetPomodoros(t)

		req, err := http.NewRequest("GET", testUrl+"/pomodoros", nil)
		if err != nil {
			t.Error("Create request failed:", err)
		}

		resp, err := testClient.Do(req)
		if err != nil {
			t.Error("Do request failed:", err)
		}

		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Error("Read response failed:", err)
		}
		if err := resp.Body.Close(); err != nil {
			t.Error("Close response failed:", err)
		}

		var body pomodorosResponse
		if err := json.Unmarshal(bytes, &body); err != nil {
			t.Error("Unmarshal json failed:", err)
		}

		if resp.StatusCode != 200 {
			t.Error("Status code should be 200, but", resp.StatusCode)
		}

		if len(body.Pomodoros) != 2 {
			t.Error("Pomodoros should have 2 pomodoro, but", len(body.Pomodoros))
		}
		pomodoroRecord1 := body.Pomodoros[0]
		if pomodoroRecord1.Task.ActualPomodoroNumber != 2 {
			t.Error("Task1's ActualPomodoroNumber should be 2, but", pomodoroRecord1.Task.ActualPomodoroNumber)
		}
		if pomodoroRecord1.ID != 1 {
			t.Error("ID should be 1, but", pomodoroRecord1.ID)
		}
		if pomodoroRecord1.CompletedAt != "2021-08-31T01:02:03Z" {
			t.Error("CompletedAt should be 2021-08-31T01:02:03Z")
		}
		if pomodoroRecord1.CreatedAt == "0001-01-01T00:00:00Z" {
			t.Error("CreatedAt should not be 0001-01-01T00:00:00Z")
		}
		pomodoroRecord2 := body.Pomodoros[1]
		if pomodoroRecord2.ID != 2 {
			t.Error("ID should be 2, but", pomodoroRecord2.ID)
		}
		if pomodoroRecord2.CompletedAt != "2021-09-01T06:07:08Z" {
			t.Error("CompletedAt should be 2021-09-01T06:07:08Z")
		}
		if pomodoroRecord2.CreatedAt == "0001-01-01T00:00:00Z" {
			t.Error("CreatedAt should not be 0001-01-01T00:00:00Z")
		}

		shutdownTestDB(t)
	})
	t.Run("2021年8月31日に実行したポモドーロの記録を一覧取得する", func(t *testing.T) {
		setupTestDB(t)
		setupTestGetPomodoros(t)

		req, err := http.NewRequest("GET", testUrl+"/pomodoros", nil)
		if err != nil {
			t.Error("Create request failed:", err)
		}

		params := req.URL.Query()
		params.Add("completed-on", "2021-08-31T00:00:00Z")
		req.URL.RawQuery = params.Encode()

		resp, err := testClient.Do(req)
		if err != nil {
			t.Error("Do request failed:", err)
		}

		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Error("Read response failed:", err)
		}
		if err := resp.Body.Close(); err != nil {
			t.Error("Close response failed:", err)
		}

		var body pomodorosResponse
		if err := json.Unmarshal(bytes, &body); err != nil {
			t.Error("Unmarshal json failed:", err)
		}

		if resp.StatusCode != 200 {
			t.Error("Status code should be 200, but", resp.StatusCode)
		}

		if len(body.Pomodoros) != 1 {
			t.Error("Pomodoros should have 1 pomodoro, but", len(body.Pomodoros))
		}
		pomodoroRecord1 := body.Pomodoros[0]
		if pomodoroRecord1.Task.ActualPomodoroNumber != 2 {
			t.Error("Task1's ActualPomodoroNumber should be 2, but", pomodoroRecord1.Task.ActualPomodoroNumber)
		}
		if pomodoroRecord1.ID != 1 {
			t.Error("ID should be 1, but", pomodoroRecord1.ID)
		}
		if pomodoroRecord1.CompletedAt != "2021-08-31T01:02:03Z" {
			t.Error("CompletedAt should be 2021-08-31T01:02:03Z")
		}
		if pomodoroRecord1.CreatedAt == "0001-01-01T00:00:00Z" {
			t.Error("CreatedAt should not be 0001-01-01T00:00:00Z")
		}

		shutdownTestDB(t)
	})
}

func TestGetRestCount(t *testing.T) {
	t.Run("次の15分休憩までのカウントを取得する", func(t *testing.T) {
		setupTestDB(t)

		req, err := http.NewRequest("GET", testUrl+"/pomodoros/rest-count", nil)
		if err != nil {
			t.Error("Create request failed:", err)
		}

		resp, err := testClient.Do(req)
		if err != nil {
			t.Error("Do request failed:", err)
		}

		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Error("Read response failed:", err)
		}
		if err := resp.Body.Close(); err != nil {
			t.Error("Close response failed:", err)
		}

		var body restCountResponse
		if err := json.Unmarshal(bytes, &body); err != nil {
			t.Error("Unmarshal json failed:", err)
		}

		if resp.StatusCode != 200 {
			t.Error("Status code should be 200, but", resp.StatusCode)
		}
		if body.RestCount != 4 {
			t.Error("restCount should be 4, but", body.RestCount)
		}

		shutdownTestDB(t)
	})
}

func BenchmarkPostPomodoro(b *testing.B) {
	setupTestDB(b)
	setupTestPostPomodoro(b)
	defer shutdownTestDB(b)

	for i := 0; i < b.N; i++ {
		reqBody := strings.NewReader(`{ "taskID": 1 }`)
		req, _ := http.NewRequest("POST", testUrl+"/pomodoros", reqBody)

		_, _ = testClient.Do(req)
	}
}

func BenchmarkGetPomodoros(b *testing.B) {
	setupTestDB(b)
	setupTestGetPomodoros(b)
	defer shutdownTestDB(b)

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", testUrl+"/pomodoros", nil)

		_, _ = testClient.Do(req)
	}
}

func BenchmarkGetRestCount(b *testing.B) {
	setupTestDB(b)
	defer shutdownTestDB(b)

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", testUrl+"/pomodoros/rest-count", nil)

		_, _ = testClient.Do(req)
	}
}

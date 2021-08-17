package tomeit

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func setupTestPostPomodoroRecord(tb testing.TB) {
	const createTask1 = `INSERT INTO tasks (user_id, title, expected_pomodoro_number, due_on, is_completed) VALUES (1, 'タスク1', 0, '2018-12-31', false)`

	if _, err := testDB.Exec(createTask1); err != nil {
		tb.Fatal("setupTestPostPomodoroLog failed:", err)
	}
}

func TestPostPomodoroRecord(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		setupTestDB(t)
		setupTestPostPomodoroRecord(t)

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
		if task.CreatedAt == "" {
			t.Error("CreatedAt does not exist")
		}
		if task.UpdatedAt == "" {
			t.Error("UpdatedAt does not exist")
		}
		if body.CreatedAt == "" {
			t.Error("CreatedAt does not exist")
		}

		shutdownTestDB(t)
	})
}

//
//func setupTestGetPomodoroRecords(tb testing.TB) {
//	const createTask1 = `INSERT INTO tasks (user_id, title, expectedPomodoroNumber, dueOn, is_done) VALUES (1, 'タスク1', 0, '2021-01-01', false)`
//	const createPomodoroLog1 = `INSERT INTO pomodoro_logs (user_id, task_id) VALUES (1, 1)`
//	const createPomodoroLog2 = `INSERT INTO pomodoro_logs (user_id, task_id) VALUES (1, 1)`
//
//	if _, err := testDB.Exec(createTask1); err != nil {
//		tb.Fatal("setupTestPostPomodoroLog failed:", err)
//	}
//	if _, err := testDB.Exec(createPomodoroLog1); err != nil {
//		tb.Fatal("setupTestPostPomodoroLog failed:", err)
//	}
//	time.Sleep(time.Second)
//	if _, err := testDB.Exec(createPomodoroLog2); err != nil {
//		tb.Fatal("setupTestPostPomodoroLog failed:", err)
//	}
//}
//
//func TestGetPomodoroRecords(t *testing.T) {
//	t.Run("success", func(t *testing.T) {
//		setupTestDB()
//		setupTestGetPomodoroRecords(t)
//
//		req, err := http.NewRequest("GET", testUrl+"/pomodoros/records", nil)
//		if err != nil {
//			t.Error("Create request failed:", err)
//		}
//
//		resp, err := testClient.Do(req)
//		if err != nil {
//			t.Error("Do request failed:", err)
//		}
//
//		bytes, err := io.ReadAll(resp.Body)
//		if err != nil {
//			t.Error("Read response failed:", err)
//		}
//		if err := resp.Body.Close(); err != nil {
//			t.Error("Close response failed:", err)
//		}
//
//		var body pomodorosResponse
//		if err := json.Unmarshal(bytes, &body); err != nil {
//			t.Error("Unmarshal json failed:", err)
//		}
//
//		if resp.StatusCode != 200 {
//			t.Error("Status code should be 200, but", resp.StatusCode)
//		}
//
//		if len(body.Pomodoros) != 2 {
//			t.Error("Pomodoros should have 2 pomodoro, but", len(body.Pomodoros))
//		}
//
//		pomodoroRecord1 := body.Pomodoros[0]
//		if pomodoroRecord1.ID != 1 {
//			t.Error("ID should be 1, but", pomodoroRecord1.ID)
//		}
//		if pomodoroRecord1.CreatedAt == "" {
//			t.Error("CreatedAt does not exist")
//		}
//
//		pomodoroRecord2 := body.Pomodoros[1]
//		if pomodoroRecord2.ID != 2 {
//			t.Error("ID should be 2, but", pomodoroRecord2.ID)
//		}
//		if pomodoroRecord2.CreatedAt == "" {
//			t.Error("CreatedAt does not exist")
//		}
//
//		shutdownTestDB()
//	})
//}
//
//func TestGetRestCount(t *testing.T) {
//	t.Run("success", func(t *testing.T) {
//		setupTestDB()
//
//		req, err := http.NewRequest("GET", testUrl+"/pomodoros/rest-count", nil)
//		if err != nil {
//			t.Error("Create request failed:", err)
//		}
//
//		resp, err := testClient.Do(req)
//		if err != nil {
//			t.Error("Do request failed:", err)
//		}
//
//		bytes, err := io.ReadAll(resp.Body)
//		if err != nil {
//			t.Error("Read response failed:", err)
//		}
//		if err := resp.Body.Close(); err != nil {
//			t.Error("Close response failed:", err)
//		}
//
//		var body restCountResponse
//		if err := json.Unmarshal(bytes, &body); err != nil {
//			t.Error("Unmarshal json failed:", err)
//		}
//
//		if resp.StatusCode != 200 {
//			t.Error("Status code should be 200, but", resp.StatusCode)
//		}
//
//		if body.RestCount != 4 {
//			t.Error("nextRestCount should be 4, but", body.RestCount)
//		}
//
//		shutdownTestDB()
//	})
//}

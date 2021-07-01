package tomeit

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func setupTestPostPomodoroLog(tb testing.TB) {
	const createTask1 = `INSERT INTO tasks (user_id, name, priority, deadline, is_done) VALUES (1, 'タスク1', 0, '2021-01-01', false)`

	if _, err := testDB.Exec(createTask1); err != nil {
		tb.Fatal("setupTestPostPomodoroLog failed:", err)
	}
}

func TestPostPomodoroLog(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		setupTestDB()
		setupTestPostPomodoroLog(t)

		reqBody := strings.NewReader(`{"taskID": 1 }`)
		req, err := http.NewRequest("POST", testUrl+"/pomodoros/logs", reqBody)
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

		var body pomodoroLogResponse
		if err := json.Unmarshal(bytes, &body); err != nil {
			t.Error("Unmarshal json failed:", err)
		}

		if resp.StatusCode != 201 {
			t.Error("Status code should be 201, but", resp.StatusCode)
		}

		if body.ID != 1 {
			t.Error("ID should be 1, but", body.ID)
		}
		if body.CreatedAt == "" {
			t.Error("CreatedAt does not exist")
		}

		task := body.Task

		if task.ID != 1 {
			t.Error("ID should be 1, but", task.ID)
		}
		if task.Name != "タスク1" {
			t.Error("Name should be タスク2, but", task.Name)
		}
		if task.Priority != 0 {
			t.Error("Priority should be 1, but", task.Priority)
		}
		if task.Deadline != "2021-01-01" {
			t.Error("Deadline should be 2021-07-01, but", task.Deadline)
		}
		if task.IsDone != false {
			t.Error("IsDone should be false, but", task.IsDone)
		}
		if task.PomodoroCount != 1 {
			t.Error("PomodoroCount should be 1, but", task.PomodoroCount)
		}
		if task.CreatedAt == "" {
			t.Error("CreatedAt does not exist")
		}
		if task.UpdatedAt == "" {
			t.Error("UpdatedAt does not exist")
		}

		shutdownTestDB()
	})
}

func setupTestGetPomodoroLogs(tb testing.TB) {
	const createTask1 = `INSERT INTO tasks (user_id, name, priority, deadline, is_done) VALUES (1, 'タスク1', 0, '2021-01-01', false)`
	const createPomodoroLog1 = `INSERT INTO pomodoro_logs (user_id, task_id) VALUES (1, 1)`
	const createPomodoroLog2 = `INSERT INTO pomodoro_logs (user_id, task_id) VALUES (1, 1)`

	if _, err := testDB.Exec(createTask1); err != nil {
		tb.Fatal("setupTestPostPomodoroLog failed:", err)
	}
	if _, err := testDB.Exec(createPomodoroLog1); err != nil {
		tb.Fatal("setupTestPostPomodoroLog failed:", err)
	}
	time.Sleep(time.Second)
	if _, err := testDB.Exec(createPomodoroLog2); err != nil {
		tb.Fatal("setupTestPostPomodoroLog failed:", err)
	}
}

func TestGetPomodoroLogs(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		setupTestDB()
		setupTestGetPomodoroLogs(t)

		req, err := http.NewRequest("GET", testUrl+"/pomodoros/logs", nil)
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

		var body pomodoroLogsResponse
		if err := json.Unmarshal(bytes, &body); err != nil {
			t.Error("Unmarshal json failed:", err)
		}

		if resp.StatusCode != 200 {
			t.Error("Status code should be 200, but", resp.StatusCode)
		}

		if len(body.PomodoroLogs) != 2 {
			t.Error("PomodoroLogs should have 2 pomodoroLog, but", len(body.PomodoroLogs))
		}

		pomodoroLog2 := body.PomodoroLogs[0]
		if pomodoroLog2.ID != 2 {
			t.Error("ID should be 2, but", pomodoroLog2.ID)
		}
		if pomodoroLog2.CreatedAt == "" {
			t.Error("CreatedAt does not exist")
		}

		pomodoroLog1 := body.PomodoroLogs[1]
		if pomodoroLog1.ID != 1 {
			t.Error("ID should be 1, but", pomodoroLog1.ID)
		}
		if pomodoroLog1.CreatedAt == "" {
			t.Error("CreatedAt does not exist")
		}

		shutdownTestDB()
	})
}

func TestGetRestCount(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		setupTestDB()

		req, err := http.NewRequest("GET", testUrl+"/pomodoros/rest/count", nil)
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

		if body.CountToNextRest != 4 {
			t.Error("CountToNextRest should be 4, but", body.CountToNextRest)
		}

		shutdownTestDB()
	})
}

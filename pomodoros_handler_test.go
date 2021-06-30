package tomeit

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestPostPomodoroLog(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		reqBody := strings.NewReader(`{"taskID": 2 }`)
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
		if body.ID <= 0 {
			t.Error("Id should be positive number, but", body.ID)
		}
		if body.CreatedAt == "" {
			t.Error("CreatedAt does not exist")
		}

		task := body.Task

		if task.Name != "タスク2" {
			t.Error("Name should be タスク2, but", task.Name)
		}
		if task.Priority != 1 {
			t.Error("Priority should be 1, but", task.Priority)
		}
		if task.Deadline != "2021-07-01" {
			t.Error("Deadline should be 2021-07-01, but", task.Deadline)
		}
		if task.IsDone != false {
			t.Error("IsDone should be false, but", task.IsDone)
		}
		if task.PomodoroCount != 0 {
			t.Error("PomodoroCount should be 0, but", task.PomodoroCount)
		}
		if task.CreatedAt == "" {
			t.Error("CreatedAt does not exist")
		}
		if task.UpdatedAt == "" {
			t.Error("UpdatedAt does not exist")
		}
	})
}

func TestGetPomodoroLogs(t *testing.T) {
	t.Run("success", func(t *testing.T) {
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

		if len(body.PomodoroLogs) != 1 {
			t.Error("Tasks should have five task, but", len(body.PomodoroLogs))
		}
	})
}

package tomeit

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestPostTask(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		reqBody := strings.NewReader(`{"name": "タスク1", "priority": 2, "deadline": "2021-06-30"}`)
		req, err := http.NewRequest("POST", testUrl+"/tasks", reqBody)
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

		var body taskResponse
		if err := json.Unmarshal(bytes, &body); err != nil {
			t.Error("Unmarshal json failed:", err)
		}

		if resp.StatusCode != 201 {
			t.Error("Status code should be 201, but", resp.StatusCode)
		}
		if body.ID <= 0 {
			t.Error("Id should be positive number, but", body.ID)
		}
		if body.Name != "タスク1" {
			t.Error("Name should be タスク1, but", body.Name)
		}
		if body.Priority != 2 {
			t.Error("Priority should be 2, but", body.Priority)
		}
		if body.Deadline != "2021-06-30" {
			t.Error("Deadline should be 2021-06-30, but", body.Deadline)
		}
		if body.IsDone != false {
			t.Error("IsDone should be false, but", body.IsDone)
		}
		if body.PomodoroCount != 0 {
			t.Error("PomodoroCount should be 0, but", body.PomodoroCount)
		}
		if body.CreatedAt == "" {
			t.Error("CreatedAt does not exist")
		}
		if body.UpdatedAt == "" {
			t.Error("UpdatedAt does not exist")
		}
	})
}

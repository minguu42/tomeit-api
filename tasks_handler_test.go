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
		setupTestDB()

		reqBody := strings.NewReader(`{"name": "新しいタスク", "priority": 1, "deadline": "2021-01-01"}`)
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
		if body.Name != "新しいタスク" {
			t.Error("Name should be 新しいタスク, but", body.Name)
		}
		if body.Priority != 1 {
			t.Error("Priority should be 1, but", body.Priority)
		}
		if body.Deadline != "2021-01-01" {
			t.Error("Deadline should be 2021-01-01, but", body.Deadline)
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

		shutdownTestDB()
	})
}

//func TestGetTasks(t *testing.T) {
//	t.Run("success", func(t *testing.T) {
//		req, err := http.NewRequest("GET", testUrl+"/tasks", nil)
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
//		var body tasksResponse
//		if err := json.Unmarshal(bytes, &body); err != nil {
//			t.Error("Unmarshal json failed:", err)
//		}
//
//		if resp.StatusCode != 200 {
//			t.Error("Status code should be 200, but", resp.StatusCode)
//		}
//
//		if len(body.Tasks) != 5 {
//			t.Error("Tasks should have five task4, but", len(body.Tasks))
//		}
//
//		task4 := body.Tasks[1]
//		if task4.ID != 4 {
//			t.Error("Id should be 4, but", task4.ID)
//		}
//		if task4.Name != "タスク4" {
//			t.Error("Name should be タスク4, but", task4.Name)
//		}
//		if task4.Priority != 3 {
//			t.Error("Priority should be 3, but", task4.Priority)
//		}
//		if task4.Deadline != "2021-01-04" {
//			t.Error("Deadline should be 2021-01-04, but", task4.Deadline)
//		}
//		if task4.IsDone != true {
//			t.Error("IsDone should be true, but", task4.IsDone)
//		}
//		if task4.PomodoroCount != 0 {
//			t.Error("PomodoroCount should be 0, but", task4.PomodoroCount)
//		}
//		if task4.CreatedAt == "" {
//			t.Error("CreatedAt does not exist")
//		}
//		if task4.UpdatedAt == "" {
//			t.Error("UpdatedAt does not exist")
//		}
//
//		task1 := body.Tasks[4]
//		if task1.ID != 1 {
//			t.Error("Id should be 1, but", task1.ID)
//		}
//		if task1.Name != "タスク1" {
//			t.Error("Name should be タスク1, but", task1.Name)
//		}
//		if task1.Priority != 0 {
//			t.Error("Priority should be 0, but", task1.Priority)
//		}
//		if task1.Deadline != "2021-01-01" {
//			t.Error("Deadline should be 2021-01-01, but", task1.Deadline)
//		}
//		if task1.IsDone != false {
//			t.Error("IsDone should be false, but", task1.IsDone)
//		}
//		if task1.PomodoroCount != 0 {
//			t.Error("PomodoroCount should be 0, but", task1.PomodoroCount)
//		}
//		if task1.CreatedAt == "" {
//			t.Error("CreatedAt does not exist")
//		}
//		if task1.UpdatedAt == "" {
//			t.Error("UpdatedAt does not exist")
//		}
//	})
//}
//
//func TestGetTasksDone(t *testing.T) {
//	t.Run("success", func(t *testing.T) {
//		req, err := http.NewRequest("GET", testUrl+"/tasks/done", nil)
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
//		var body tasksResponse
//		if err := json.Unmarshal(bytes, &body); err != nil {
//			t.Error("Unmarshal json failed:", err)
//		}
//
//		if resp.StatusCode != 200 {
//			t.Error("Status code should be 200, but", resp.StatusCode)
//		}
//
//		if len(body.Tasks) != 2 {
//			t.Error("Tasks should have two task, but", len(body.Tasks))
//		}
//	})
//}
//
//func TestPutTaskDone(t *testing.T) {
//	t.Run("success", func(t *testing.T) {
//		req, err := http.NewRequest("PUT", testUrl+"/tasks/done/1", nil)
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
//		var body taskResponse
//		if err := json.Unmarshal(bytes, &body); err != nil {
//			t.Error("Unmarshal json failed:", err)
//		}
//
//		if resp.StatusCode != 200 {
//			t.Error("Status code should be 201, but", resp.StatusCode)
//		}
//		if body.ID != 1 {
//			t.Error("Id should be 1, but", body.ID)
//		}
//		if body.Name != "タスク1" {
//			t.Error("Name should be タスク1, but", body.Name)
//		}
//		if body.Priority != 0 {
//			t.Error("Priority should be 0, but", body.Priority)
//		}
//		if body.Deadline != "2021-06-30" {
//			t.Error("Deadline should be 2021-06-30, but", body.Deadline)
//		}
//		if body.IsDone != true {
//			t.Error("IsDone should be true, but", body.IsDone)
//		}
//		if body.PomodoroCount != 0 {
//			t.Error("PomodoroCount should be 0, but", body.PomodoroCount)
//		}
//		if body.CreatedAt == "" {
//			t.Error("CreatedAt does not exist")
//		}
//		if body.UpdatedAt == "" {
//			t.Error("UpdatedAt does not exist")
//		}
//	})
//}

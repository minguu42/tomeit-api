package tomeit

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPostTasks(t *testing.T) {
	t.Run("新しいタスクを作成する", func(t *testing.T) {
		setupTestDB(t)
		t.Cleanup(teardownTestDB)

		reqBody := strings.NewReader(`
{
  "title": "タスク1",
  "expectedPomodoroNum": 4,
  "dueOn": "2021-01-01T00:00:00Z"
}
`)
		resp, body := doTestRequest(t, "POST", "/tasks", reqBody, "taskResponse")
		if resp.StatusCode != 200 {
			t.Error("Status code should be 200, but", resp.StatusCode)
		}

		got, ok := body.(taskResponse)
		if !ok {
			t.Fatal("Type Assertion failed")
		}

		want := taskResponse{
			ID:                  1,
			Title:               "タスク1",
			ExpectedPomodoroNum: 4,
			ActualPomodoroNum:   0,
			DueOn:               "2021-01-01T00:00:00Z",
			IsCompleted:         false,
			CompletedOn:         "",
		}

		if diff := cmp.Diff(got, want, taskResponseCmpOpts); diff != "" {
			t.Errorf("taskResponse mismatch (-got +want):\n%s", diff)
		}
	})
	t.Run("リクエストボディに title が存在しない", func(t *testing.T) {
		setupTestDB(t)
		t.Cleanup(teardownTestDB)

		reqBody := strings.NewReader(`{
  "expectedPomodoroNum": 4,
  "dueOn": "2021-01-01T00:00:00Z"
}`)
		resp, _ := doTestRequest(t, "POST", "/tasks", reqBody, "taskResponse")
		if resp.StatusCode != 400 {
			t.Error("Status code should be 400, but", resp.StatusCode)
		}
	})
}

//func TestPostTasks(t *testing.T) {
//	t.Run("新しいタスクを作成する", func(t *testing.T) {
//		setupTestDB(t)
//
//		reqBody := strings.NewReader(`{"title": "タスク1", "expectedPomodoroNumber": 2, "dueOn": "0001-01-01T00:00:00Z"}`)
//		req, err := http.NewRequest("POST", testUrl+"/tasks", reqBody)
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
//		if resp.StatusCode != 201 {
//			t.Error("Status code should be 201, but", resp.StatusCode)
//		}
//		if body.ID != 1 {
//			t.Error("ID should be 1, but", body.ID)
//		}
//		if body.Title != "タスク1" {
//			t.Error("Title should be タスク1, but", body.Title)
//		}
//		if body.ExpectedPomodoroNumber != 2 {
//			t.Error("ExpectedPomodoroNumber should be 2, but", body.ExpectedPomodoroNumber)
//		}
//		if body.ActualPomodoroNumber != 0 {
//			t.Error("ActualPomodoroNumber should be 0, but", body.ActualPomodoroNumber)
//		}
//		if body.DueOn != "0001-01-01T00:00:00Z" {
//			t.Error("DueOn should be 0001-01-01T00:00:00Z, but", body.DueOn)
//		}
//		if body.IsCompleted != false {
//			t.Error("IsCompleted should be false, but", body.IsCompleted)
//		}
//		if body.CompletedAt != "0001-01-01T00:00:00Z" {
//			t.Error("CompletedAt should be 0001-01-01T00:00:00Z, but", body.CompletedAt)
//		}
//		if body.CreatedAt == "0001-01-01T00:00:00Z" {
//			t.Error("CreatedAt does not exist")
//		}
//		if body.UpdatedAt == "0001-01-01T00:00:00Z" {
//			t.Error("UpdatedAt does not exist")
//		}
//
//		shutdownTestDB(t)
//	})
//}
//
//func setupTestGetTasks(tb testing.TB) {
//	const createTask1 = `INSERT INTO tasks (user_id, title, expected_pomodoro_number, due_on, is_completed) VALUES (1, 'タスク1', 0, '2021-01-01', false)`
//	const createTask2 = `INSERT INTO tasks (user_id, title, expected_pomodoro_number, due_on, is_completed, completed_at) VALUES (1, 'タスク2', 1, '2021-12-31', true, '2021-08-31 12:30:00')`
//
//	if _, err := testDB.Exec(createTask1); err != nil {
//		tb.Fatal("create Task1 failed:", err)
//	}
//	if _, err := testDB.Exec(createTask2); err != nil {
//		tb.Fatal("create Task2 failed:", err)
//	}
//}
//
//func TestGetTasks(t *testing.T) {
//	t.Run("タスク一覧を取得する", func(t *testing.T) {
//		setupTestDB(t)
//		setupTestGetTasks(t)
//
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
//		if len(body.Tasks) != 2 {
//			t.Error("Tasks should have 2 tasks, but", len(body.Tasks))
//		}
//
//		task1 := body.Tasks[0]
//		if task1.ID != 1 {
//			t.Error("ID should be 1, but", task1.ID)
//		}
//		if task1.Title != "タスク1" {
//			t.Error("Title should be タスク1, but", task1.Title)
//		}
//		if task1.ExpectedPomodoroNumber != 0 {
//			t.Error("ExpectedPomodoroNumber should be 0, but", task1.ExpectedPomodoroNumber)
//		}
//		if task1.DueOn != "2021-01-01T00:00:00Z" {
//			t.Error("DueOn should be 2021-01-01T00:00:00Z, but", task1.DueOn)
//		}
//		if task1.IsCompleted != false {
//			t.Error("IsCompleted should be false, but", task1.IsCompleted)
//		}
//		if task1.CompletedAt != "0001-01-01T00:00:00Z" {
//			t.Error("CompletedAt should be 0001-01-01T00:00:00Z, but", task1.CompletedAt)
//		}
//		if task1.ActualPomodoroNumber != 0 {
//			t.Error("ActualPomodoroNumber should be 0, but", task1.ActualPomodoroNumber)
//		}
//		if task1.CreatedAt == "0001-01-01T00:00:00Z" {
//			t.Error("CreatedAt does not exist")
//		}
//		if task1.UpdatedAt == "0001-01-01T00:00:00Z" {
//			t.Error("UpdatedAt does not exist")
//		}
//
//		task2 := body.Tasks[1]
//		if task2.ID != 2 {
//			t.Error("ID should be 2, but", task2.ID)
//		}
//		if task2.Title != "タスク2" {
//			t.Error("Title should be タスク2, but", task2.Title)
//		}
//		if task2.ExpectedPomodoroNumber != 1 {
//			t.Error("ExpectedPomodoroNumber should be 1, but", task2.ExpectedPomodoroNumber)
//		}
//		if task2.DueOn != "2021-12-31T00:00:00Z" {
//			t.Error("DueOn should be 2021-12-31T00:00:00Z, but", task2.DueOn)
//		}
//		if task2.IsCompleted != true {
//			t.Error("IsCompleted should be true, but", task2.IsCompleted)
//		}
//		if task2.CompletedAt != "2021-08-31T12:30:00Z" {
//			t.Error("CompletedAt should be 2021-08-31T12:30:00Z, but", task1.CompletedAt)
//		}
//		if task2.ActualPomodoroNumber != 0 {
//			t.Error("ActualPomodoroNumber should be 0, but", task2.ActualPomodoroNumber)
//		}
//		if task2.CreatedAt == "0001-01-01T00:00:00Z" {
//			t.Error("CreatedAt does not exist")
//		}
//		if task2.UpdatedAt == "0001-01-01T00:00:00Z" {
//			t.Error("UpdatedAt does not exist")
//		}
//
//		shutdownTestDB(t)
//	})
//	t.Run("完了済みのタスク一覧を取得する", func(t *testing.T) {
//		setupTestDB(t)
//		setupTestGetTasks(t)
//
//		req, err := http.NewRequest("GET", testUrl+"/tasks", nil)
//		if err != nil {
//			t.Error("Create request failed:", err)
//		}
//
//		params := req.URL.Query()
//		params.Add("is-completed", "true")
//		req.URL.RawQuery = params.Encode()
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
//		if len(body.Tasks) != 1 {
//			t.Error("Tasks should have 1 tasks, but", len(body.Tasks))
//		}
//
//		task2 := body.Tasks[0]
//		if task2.ID != 2 {
//			t.Error("ID should be 2, but", task2.ID)
//		}
//		if task2.Title != "タスク2" {
//			t.Error("Title should be タスク2, but", task2.Title)
//		}
//		if task2.ExpectedPomodoroNumber != 1 {
//			t.Error("ExpectedPomodoroNumber should be 1, but", task2.ExpectedPomodoroNumber)
//		}
//		if task2.DueOn != "2021-12-31T00:00:00Z" {
//			t.Error("DueOn should be 2021-12-31T00:00:00Z, but", task2.DueOn)
//		}
//		if task2.IsCompleted != true {
//			t.Error("IsCompleted should be true, but", task2.IsCompleted)
//		}
//		if task2.CompletedAt != "2021-08-31T12:30:00Z" {
//			t.Error("CompletedAt should be 2021-08-31T12:30:00Z, but", task2.CompletedAt)
//		}
//		if task2.ActualPomodoroNumber != 0 {
//			t.Error("ActualPomodoroNumber should be 0, but", task2.ActualPomodoroNumber)
//		}
//		if task2.CreatedAt == "0001-01-01T00:00:00Z" {
//			t.Error("CreatedAt does not exist")
//		}
//		if task2.UpdatedAt == "0001-01-01T00:00:00Z" {
//			t.Error("UpdatedAt does not exist")
//		}
//	})
//	t.Run("2021年8月31日に完了したタスク一覧を取得する", func(t *testing.T) {
//		setupTestDB(t)
//		setupTestGetTasks(t)
//
//		req, err := http.NewRequest("GET", testUrl+"/tasks", nil)
//		if err != nil {
//			t.Error("Create request failed:", err)
//		}
//
//		params := req.URL.Query()
//		params.Add("completed-on", "2021-08-31T00:00:00Z")
//		req.URL.RawQuery = params.Encode()
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
//		if len(body.Tasks) != 1 {
//			t.Error("Tasks should have 1 tasks, but", len(body.Tasks))
//		}
//
//		task2 := body.Tasks[0]
//		if task2.ID != 2 {
//			t.Error("ID should be 2, but", task2.ID)
//		}
//		if task2.Title != "タスク2" {
//			t.Error("Title should be タスク2, but", task2.Title)
//		}
//		if task2.ExpectedPomodoroNumber != 1 {
//			t.Error("ExpectedPomodoroNumber should be 1, but", task2.ExpectedPomodoroNumber)
//		}
//		if task2.ActualPomodoroNumber != 0 {
//			t.Error("ActualPomodoroNumber should be 0, but", task2.ActualPomodoroNumber)
//		}
//		if task2.DueOn != "2021-12-31T00:00:00Z" {
//			t.Error("DueOn should be 2021-12-31T00:00:00Z, but", task2.DueOn)
//		}
//		if task2.IsCompleted != true {
//			t.Error("IsCompleted should be true, but", task2.IsCompleted)
//		}
//		if task2.CompletedAt != "2021-08-31T12:30:00Z" {
//			t.Error("CompletedAt should be 2021-08-31T12:30:00Z, but", task2.CompletedAt)
//		}
//		if task2.CreatedAt == "0001-01-01T00:00:00Z" {
//			t.Error("CreatedAt does not exist")
//		}
//		if task2.UpdatedAt == "0001-01-01T00:00:00Z" {
//			t.Error("UpdatedAt does not exist")
//		}
//	})
//}
//
//func setupTestPatchTask(tb testing.TB) {
//	const createTask1 = `INSERT INTO tasks (user_id, title, expected_pomodoro_number, is_completed) VALUES (1, 'タスク1', 0, false)`
//	const createTask2 = `INSERT INTO tasks (user_id, title, expected_pomodoro_number, is_completed) VALUES (1, 'タスク2', 3, true)`
//
//	if _, err := testDB.Exec(createTask1); err != nil {
//		tb.Fatal("createTask1 failed:", err)
//	}
//	if _, err := testDB.Exec(createTask2); err != nil {
//		tb.Fatal("createTask2 failed:", err)
//	}
//}
//
//func TestPatchTask(t *testing.T) {
//	t.Run("タスク1の isCompleted の値を true に変更する", func(t *testing.T) {
//		setupTestDB(t)
//		setupTestPatchTask(t)
//
//		reqBody := strings.NewReader(`{"isCompleted": "true"}`)
//		req, err := http.NewRequest("PATCH", testUrl+"/tasks/1", reqBody)
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
//			t.Error("Status code should be 200, but", resp.StatusCode)
//		}
//		if body.ID != 1 {
//			t.Error("Id should be 1, but", body.ID)
//		}
//		if body.Title != "タスク1" {
//			t.Error("Title should be タスク1, but", body.Title)
//		}
//		if body.ExpectedPomodoroNumber != 0 {
//			t.Error("ExpectedPomodoroNumber should be 0, but", body.ExpectedPomodoroNumber)
//		}
//		if body.ActualPomodoroNumber != 0 {
//			t.Error("ActualPomodoroNumber should be 0, but", body.ActualPomodoroNumber)
//		}
//		if body.DueOn != "0001-01-01T00:00:00Z" {
//			t.Error("DueOn should be 0001-01-01, but", body.DueOn)
//		}
//		if body.IsCompleted != true {
//			t.Error("IsCompleted should be true, but", body.IsCompleted)
//		}
//		if body.CompletedAt == "" {
//			t.Error("CompletedAt does not exist")
//		}
//		if body.CreatedAt == "" {
//			t.Error("CreatedAt does not exist")
//		}
//		if body.UpdatedAt == "" {
//			t.Error("UpdatedAt does not exist")
//		}
//
//		shutdownTestDB(t)
//	})
//	t.Run("タスク2の isCompleted の値を false に変更する", func(t *testing.T) {
//		setupTestDB(t)
//		setupTestPatchTask(t)
//
//		reqBody := strings.NewReader(`{"isCompleted": "false"}`)
//		req, err := http.NewRequest("PATCH", testUrl+"/tasks/2", reqBody)
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
//			t.Error("Status code should be 200, but", resp.StatusCode)
//		}
//		if body.ID != 2 {
//			t.Error("Id should be 1, but", body.ID)
//		}
//		if body.Title != "タスク2" {
//			t.Error("Title should be タスク1, but", body.Title)
//		}
//		if body.ExpectedPomodoroNumber != 3 {
//			t.Error("ExpectedPomodoroNumber should be 0, but", body.ExpectedPomodoroNumber)
//		}
//		if body.ActualPomodoroNumber != 0 {
//			t.Error("ActualPomodoroNumber should be 0, but", body.ActualPomodoroNumber)
//		}
//		if body.DueOn != "0001-01-01T00:00:00Z" {
//			t.Error("DueOn should be 0001-01-01, but", body.DueOn)
//		}
//		if body.IsCompleted != false {
//			t.Error("IsCompleted should be false, but", body.IsCompleted)
//		}
//		if body.CompletedAt == "" {
//			t.Error("CompletedAt does not exist")
//		}
//		if body.CreatedAt == "" {
//			t.Error("CreatedAt does not exist")
//		}
//		if body.UpdatedAt == "" {
//			t.Error("UpdatedAt does not exist")
//		}
//
//		shutdownTestDB(t)
//	})
//}
//
//func BenchmarkPostTasks(b *testing.B) {
//	setupTestDB(b)
//
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		reqBody := strings.NewReader(`{"title": "タスク", "expectedPomodoroNumber": 0, "dueOn": "0001-01-01T00:00:00Z"}`)
//		req, _ := http.NewRequest("POST", testUrl+"/tasks", reqBody)
//
//		_, _ = testClient.Do(req)
//	}
//
//	shutdownTestDB(b)
//}
//
//func BenchmarkGetTasks(b *testing.B) {
//	setupTestDB(b)
//	setupTestGetTasks(b)
//
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		req, _ := http.NewRequest("GET", testUrl+"/tasks", nil)
//
//		_, _ = testClient.Do(req)
//	}
//
//	shutdownTestDB(b)
//}
//
//func BenchmarkPatchTask(b *testing.B) {
//	setupTestDB(b)
//	setupTestPatchTask(b)
//
//	for i := 0; i < b.N; i++ {
//		reqBody := strings.NewReader(`{"isCompleted": "true"}`)
//		req, _ := http.NewRequest("PATCH", testUrl+"/tasks/1", reqBody)
//
//		_, _ = testClient.Do(req)
//	}
//
//	shutdownTestDB(b)
//}

package tomeit

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPostTasks(t *testing.T) {
	setupTestDB(t)
	t.Cleanup(teardownTestDB)
	t.Run("新しいタスクを作成する", func(t *testing.T) {
		reqBody := strings.NewReader(`
{
  "title": "タスク1",
  "expectedPomodoroNum": 4,
  "dueOn": "2021-01-01T00:00:00Z"
}
`)
		resp, body := doTestRequest(t, "POST", "/tasks", nil, reqBody, "taskResponse")

		checkStatusCode(t, resp, 201)

		l := resp.Header.Get("Location")
		if l != testUrl+"/tasks/1" {
			t.Errorf("Location should be %v, but %v", testUrl+"/v0/tasks/1", l)
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
		}
		if diff := cmp.Diff(got, want, taskResponseCmpOpts); diff != "" {
			t.Errorf("taskResponse mismatch (-got +want):\n%s", diff)
		}
	})
	t.Run("リクエストボディに title が存在しない", func(t *testing.T) {
		reqBody := strings.NewReader(`{
  "expectedPomodoroNum": 4,
  "dueOn": "2021-01-01T00:00:00Z"
}`)
		resp, _ := doTestRequest(t, "POST", "/tasks", nil, reqBody, "taskResponse")

		checkStatusCode(t, resp, 400)
	})
}

func setupTestTasks() {
	const createTask1 = `INSERT INTO tasks (user_id, title, expected_pomodoro_num, due_on, is_completed) VALUES (1, 'タスク1', 0, '2021-01-01', false)`
	const createTask2 = `INSERT INTO tasks (user_id, title, expected_pomodoro_num, due_on, is_completed, completed_on) VALUES (1, 'タスク2', 2, '2021-12-31', true, '2021-08-31 12:30:00')`

	testDB.Exec(createTask1)
	testDB.Exec(createTask2)
}

func TestGetTasks(t *testing.T) {
	setupTestDB(t)
	setupTestTasks()
	t.Cleanup(teardownTestDB)
	t.Run("タスク一覧を取得する", func(t *testing.T) {
		resp, body := doTestRequest(t, "GET", "/tasks", nil, nil, "tasksResponse")
		if resp.StatusCode != 200 {
			t.Error("Status code should be 200, but", resp.StatusCode)
		}

		got, ok := body.(tasksResponse)
		if !ok {
			t.Fatal("Type Assertion failed")
		}

		want := tasksResponse{
			Tasks: []taskResponse{
				{
					ID:                  1,
					Title:               "タスク1",
					ExpectedPomodoroNum: 0,
					ActualPomodoroNum:   0,
					DueOn:               "2021-01-01T00:00:00Z",
					IsCompleted:         false,
				},
				{
					ID:                  2,
					Title:               "タスク2",
					ExpectedPomodoroNum: 2,
					ActualPomodoroNum:   0,
					DueOn:               "2021-12-31T00:00:00Z",
					IsCompleted:         true,
				},
			},
		}

		if diff := cmp.Diff(got.Tasks[0], want.Tasks[0], taskResponseCmpOpts); diff != "" {
			t.Errorf("tasksResponse mismatch (-got +want):\n%s", diff)
		}
		if diff := cmp.Diff(got.Tasks[1], want.Tasks[1], taskResponseCmpOpts); diff != "" {
			t.Errorf("tasksResponse mismatch (-got +want):\n%s", diff)
		}
	})
	t.Run("完了済みタスク一覧を取得する", func(t *testing.T) {
		params := map[string]string{
			"isCompleted": "true",
		}
		resp, body := doTestRequest(t, "GET", "/tasks", &params, nil, "tasksResponse")

		if resp.StatusCode != 200 {
			t.Error("Status code should be 200, but", resp.StatusCode)
		}

		got, ok := body.(tasksResponse)
		if !ok {
			t.Fatal("Type Assertion failed")
		}

		want := tasksResponse{
			Tasks: []taskResponse{
				{
					ID:                  2,
					Title:               "タスク2",
					ExpectedPomodoroNum: 2,
					ActualPomodoroNum:   0,
					DueOn:               "2021-12-31T00:00:00Z",
					IsCompleted:         true,
				},
			},
		}

		if diff := cmp.Diff(got.Tasks[0], want.Tasks[0], taskResponseCmpOpts); diff != "" {
			t.Errorf("tasksResponse mismatch (-got +want):\n%s", diff)
		}
	})
	t.Run("ある日付に完了したタスク一覧を取得する", func(t *testing.T) {
		params := map[string]string{
			"completedOn": "2021-08-31T00:00:00Z",
		}
		resp, body := doTestRequest(t, "GET", "/tasks", &params, nil, "tasksResponse")

		if resp.StatusCode != 200 {
			t.Error("Status code should be 200, but", resp.StatusCode)
		}

		got, ok := body.(tasksResponse)
		if !ok {
			t.Fatal("Type Assertion failed")
		}

		want := tasksResponse{
			Tasks: []taskResponse{
				{
					ID:                  2,
					Title:               "タスク2",
					ExpectedPomodoroNum: 2,
					ActualPomodoroNum:   0,
					DueOn:               "2021-12-31T00:00:00Z",
					IsCompleted:         true,
				},
			},
		}

		if diff := cmp.Diff(got.Tasks[0], want.Tasks[0], taskResponseCmpOpts); diff != "" {
			t.Errorf("tasksResponse mismatch (-got +want):\n%s", diff)
		}
	})
}

func TestPatchTask(t *testing.T) {
	setupTestDB(t)
	setupTestTasks()
	t.Cleanup(teardownTestDB)
	t.Run("タスク1の isCompleted の値を true に変更する", func(t *testing.T) {
		reqBody := strings.NewReader(`{"isCompleted": "true"}`)
		resp, body := doTestRequest(t, "PATCH", "/tasks/1", nil, reqBody, "taskResponse")

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
			ExpectedPomodoroNum: 0,
			ActualPomodoroNum:   0,
			DueOn:               "2021-01-01T00:00:00Z",
			IsCompleted:         true,
		}

		if diff := cmp.Diff(got, want, taskResponseCmpOpts); diff != "" {
			t.Errorf("tasksResponse mismatch (-got +want):\n%s", diff)
		}
	})
	t.Run("タスク2の isCompleted の値を false に変更する", func(t *testing.T) {
		reqBody := strings.NewReader(`{"isCompleted": "false"}`)
		resp, body := doTestRequest(t, "PATCH", "/tasks/2", nil, reqBody, "taskResponse")

		if resp.StatusCode != 200 {
			t.Error("Status code should be 200, but", resp.StatusCode)
		}

		got, ok := body.(taskResponse)
		if !ok {
			t.Fatal("Type Assertion failed")
		}

		want := taskResponse{
			ID:                  2,
			Title:               "タスク2",
			ExpectedPomodoroNum: 2,
			ActualPomodoroNum:   0,
			DueOn:               "2021-12-31T00:00:00Z",
			IsCompleted:         false,
		}

		if diff := cmp.Diff(got, want, taskResponseCmpOpts); diff != "" {
			t.Errorf("tasksResponse mismatch (-got +want):\n%s", diff)
		}
	})
}

func TestDeleteTask(t *testing.T) {
	setupTestDB(t)
	setupTestTasks()
	t.Cleanup(teardownTestDB)
	t.Run("タスク1を削除する", func(t *testing.T) {
		resp, _ := doTestRequest(t, "DELETE", "/tasks/1", nil, nil, "")

		if resp.StatusCode != 204 {
			t.Error("Status code should be 204, but", resp.StatusCode)
		}
	})
}

func BenchmarkPostTasks(b *testing.B) {
	setupTestDB(b)
	b.Cleanup(teardownTestDB)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		body := strings.NewReader(`{"title": "タスク", "expectedPomodoroNum": 0, "dueOn": ""}`)
		_, _ = doTestRequest(b, "POST", "/tasks", nil, body, "taskResponse")
	}
}

func BenchmarkGetTasks(b *testing.B) {
	setupTestDB(b)
	setupTestTasks()
	b.Cleanup(teardownTestDB)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = doTestRequest(b, "GET", "/tasks", nil, nil, "tasksResponse")
	}
}

func BenchmarkPatchTask(b *testing.B) {
	setupTestDB(b)
	setupTestTasks()
	b.Cleanup(teardownTestDB)

	for i := 0; i < b.N; i++ {
		body := strings.NewReader(`{"isCompleted": "true"}`)
		_, _ = doTestRequest(b, "PATCH", "/tasks/1", nil, body, "taskResponse")
	}
}

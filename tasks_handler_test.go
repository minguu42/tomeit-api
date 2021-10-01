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
			t.Errorf("Location should be %v, but %v", testUrl+"/tasks/1", l)
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
			t.Errorf("postTasks response mismatch (-got +want):\n%s", diff)
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
	const createTasks = `
INSERT INTO tasks (id, user_id, title, expected_pomodoro_num, due_on, is_completed, completed_on, created_at, updated_at) VALUES 
(1, 1, 'タスク1', 0, '2021-01-01 00:00:00', false, NULL, '2021-01-01 00:00:00', '2021-01-01 00:00:00'),
(2, 1, 'タスク2', 2, '2021-12-31 00:00:00', true, '2021-08-31 12:34:56', '2021-01-01 00:00:00', '2021-08-31 12:34:56')
`
	testDB.Exec(createTasks)
}

func TestGetTasks(t *testing.T) {
	setupTestDB(t)
	setupTestTasks()
	t.Cleanup(teardownTestDB)
	t.Run("タスク一覧を取得する", func(t *testing.T) {
		resp, body := doTestRequest(t, "GET", "/tasks", nil, nil, "tasksResponse")

		checkStatusCode(t, resp, 200)

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
					CompletedOn:         "",
					CreatedAt:           "2021-01-01T00:00:00Z",
					UpdatedAt:           "2021-01-01T00:00:00Z",
				},
				{
					ID:                  2,
					Title:               "タスク2",
					ExpectedPomodoroNum: 2,
					ActualPomodoroNum:   0,
					DueOn:               "2021-12-31T00:00:00Z",
					IsCompleted:         true,
					CompletedOn:         "2021-08-31T12:34:56Z",
					CreatedAt:           "2021-01-01T00:00:00Z",
					UpdatedAt:           "2021-08-31T12:34:56Z",
				},
			},
		}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("getTasks response mismatch (-got +want):\n%s", diff)
		}
	})
	t.Run("完了済みタスク一覧を取得する", func(t *testing.T) {
		params := map[string]string{
			"isCompleted": "true",
		}
		resp, body := doTestRequest(t, "GET", "/tasks", &params, nil, "tasksResponse")

		checkStatusCode(t, resp, 200)

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
					CompletedOn:         "2021-08-31T12:34:56Z",
					CreatedAt:           "2021-01-01T00:00:00Z",
					UpdatedAt:           "2021-08-31T12:34:56Z",
				},
			},
		}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("getTasks response mismatch (-got +want):\n%s", diff)
		}
	})
	t.Run("ある日付に完了したタスク一覧を取得する", func(t *testing.T) {
		params := map[string]string{
			"completedOn": "2021-08-31T00:00:00Z",
		}
		resp, body := doTestRequest(t, "GET", "/tasks", &params, nil, "tasksResponse")

		checkStatusCode(t, resp, 200)

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
					CompletedOn:         "2021-08-31T12:34:56Z",
					CreatedAt:           "2021-01-01T00:00:00Z",
					UpdatedAt:           "2021-08-31T12:34:56Z",
				},
			},
		}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("getTasks response mismatch (-got +want):\n%s", diff)
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

		checkStatusCode(t, resp, 200)

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
			t.Errorf("patchTask response mismatch (-got +want):\n%s", diff)
		}
	})
	t.Run("タスク2の isCompleted の値を false に変更する", func(t *testing.T) {
		reqBody := strings.NewReader(`{"isCompleted": "false"}`)
		resp, body := doTestRequest(t, "PATCH", "/tasks/2", nil, reqBody, "taskResponse")

		checkStatusCode(t, resp, 200)

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
			t.Errorf("patchTask response mismatch (-got +want):\n%s", diff)
		}
	})
	t.Run("URL で不適切な taskID を指定した場合は400エラーを返す", func(t *testing.T) {
		reqBody := strings.NewReader(`{"isCompleted": "true"}`)
		resp, _ := doTestRequest(t, "PATCH", "/tasks/一", nil, reqBody, "taskResponse")

		checkStatusCode(t, resp, 400)
	})
	t.Run("存在しないタスクを指定した場合は404エラーを返す", func(t *testing.T) {
		reqBody := strings.NewReader(`{"isCompleted": "true"}`)
		resp, _ := doTestRequest(t, "PATCH", "/tasks/3", nil, reqBody, "taskResponse")

		checkStatusCode(t, resp, 404)
	})
}

func TestPutTask(t *testing.T) {
	setupTestDB(t)
	setupTestTasks()
	t.Cleanup(teardownTestDB)
	t.Run("タスク1の値を更新する", func(t *testing.T) {
		reqBody := strings.NewReader(`
{
  "title": "新タスク1",
  "expectedPomodoroNum": 2,
  "dueOn": "2021-01-02T00:00:00Z",
  "isCompleted": true
}
`)
		resp, body := doTestRequest(t, "PUT", "/tasks/1", nil, reqBody, "taskResponse")

		checkStatusCode(t, resp, 200)

		got, ok := body.(taskResponse)
		if !ok {
			t.Fatal("Type Assertion failed")
		}

		want := taskResponse{
			ID:                  1,
			Title:               "新タスク1",
			ExpectedPomodoroNum: 2,
			ActualPomodoroNum:   0,
			DueOn:               "2021-01-02T00:00:00Z",
			IsCompleted:         true,
		}
		if diff := cmp.Diff(got, want, taskResponseCmpOpts); diff != "" {
			t.Errorf("putTask response mismatch (-got +want):\n%s", diff)
		}
	})
	t.Run("存在しないタスクを指定した場合は404エラーを返す", func(t *testing.T) {
		reqBody := strings.NewReader(`
{
  "title": "新タスク1",
  "expectedPomodoroNum": 2,
  "dueOn": "2021-01-02T00:00:00Z",
  "isCompleted": true
}
`)
		resp, _ := doTestRequest(t, "PUT", "/tasks/3", nil, reqBody, "taskResponse")

		checkStatusCode(t, resp, 404)
	})
}

func TestDeleteTask(t *testing.T) {
	setupTestDB(t)
	setupTestTasks()
	t.Cleanup(teardownTestDB)
	t.Run("タスク1を削除する", func(t *testing.T) {
		resp, _ := doTestRequest(t, "DELETE", "/tasks/1", nil, nil, "")

		checkStatusCode(t, resp, 204)
	})
	t.Run("URL で不適切な taskID を指定した場合は400エラーを返す", func(t *testing.T) {
		resp, _ := doTestRequest(t, "DELETE", "/tasks/一", nil, nil, "")

		checkStatusCode(t, resp, 400)
	})
	t.Run("存在しないタスクを指定した場合は404エラーを返す", func(t *testing.T) {
		resp, _ := doTestRequest(t, "DELETE", "/tasks/1", nil, nil, "")

		checkStatusCode(t, resp, 404)
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

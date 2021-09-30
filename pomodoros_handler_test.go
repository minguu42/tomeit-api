package tomeit

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPostPomodoros(t *testing.T) {
	setupTestDB(t)
	setupTestTasks()
	t.Cleanup(teardownTestDB)
	t.Run("タスク1のポモドーロを記録する", func(t *testing.T) {
		reqBody := strings.NewReader(`{"taskID": 1}`)
		resp, body := doTestRequest(t, "POST", "/pomodoros", nil, reqBody, "pomodoroResponse")

		checkStatusCode(t, resp, 201)

		l := resp.Header.Get("Location")
		if l != testUrl+"/pomodoros/1" {
			t.Errorf("Location should be %v, but %v", testUrl+"/pomodoros/1", l)
		}

		got, ok := body.(pomodoroResponse)
		if !ok {
			t.Fatal("Type Assertion failed")
		}

		want := pomodoroResponse{
			ID: 1,
			Task: &taskResponse{
				ID:                  1,
				Title:               "タスク1",
				ExpectedPomodoroNum: 0,
				ActualPomodoroNum:   1,
				DueOn:               "2021-01-01T00:00:00Z",
				IsCompleted:         false,
			},
		}
		if diff := cmp.Diff(got, want, pomodoroResponseCmpOpts); diff != "" {
			t.Errorf("postPomodoros response mismatch (-got +want):\n%s", diff)
		}
	})
	t.Run("リクエストボディに taskID が含まれていない", func(t *testing.T) {
		reqBody := strings.NewReader(`{}`)
		resp, _ := doTestRequest(t, "POST", "/pomodoros", nil, reqBody, "pomodoroResponse")

		checkStatusCode(t, resp, 400)
	})
}

func setupTestPomodoros() {
	setupTestTasks()
	const createPomodoros = `
INSERT INTO pomodoros (id, user_id, task_id, created_at) VALUES 
(1, 1, 1, '2021-08-31 01:02:03'),
(2, 1, 1, '2021-09-01 06:07:08')
`

	testDB.Exec(createPomodoros)
}

func TestGetPomodoros(t *testing.T) {
	setupTestDB(t)
	setupTestPomodoros()
	t.Cleanup(teardownTestDB)
	t.Run("ポモドーロ記録を一覧取得する", func(t *testing.T) {
		resp, body := doTestRequest(t, "GET", "/pomodoros", nil, nil, "pomodorosResponse")

		checkStatusCode(t, resp, 200)

		got, ok := body.(pomodorosResponse)
		if !ok {
			t.Fatal("Type Assertion failed")
		}

		if len(got.Pomodoros) != 2 {
			t.Fatal("response has 2 pomodoros")
		}

		want := pomodorosResponse{
			Pomodoros: []*pomodoroResponse{
				{
					ID: 1,
					Task: &taskResponse{
						ID:                  1,
						Title:               "タスク1",
						ExpectedPomodoroNum: 0,
						ActualPomodoroNum:   2,
						DueOn:               "2021-01-01T00:00:00Z",
						IsCompleted:         false,
						CompletedOn:         "",
						CreatedAt:           "2021-01-01T00:00:00Z",
						UpdatedAt:           "2021-01-01T00:00:00Z",
					},
					CreatedAt: "2021-08-31T01:02:03Z",
				},
				{
					ID: 2,
					Task: &taskResponse{
						ID:                  1,
						Title:               "タスク1",
						ExpectedPomodoroNum: 0,
						ActualPomodoroNum:   2,
						DueOn:               "2021-01-01T00:00:00Z",
						IsCompleted:         false,
						CompletedOn:         "",
						CreatedAt:           "2021-01-01T00:00:00Z",
						UpdatedAt:           "2021-01-01T00:00:00Z",
					},
					CreatedAt: "2021-09-01T06:07:08Z",
				},
			},
		}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("getPomodoros response mismatch (-got +want):\n%s", diff)
		}
	})
	t.Run("ある日付に作成したポモドーロ記録を取得する", func(t *testing.T) {
		params := map[string]string{
			"completedOn": "2021-08-31T00:00:00Z",
		}
		resp, body := doTestRequest(t, "GET", "/pomodoros", &params, nil, "pomodorosResponse")

		checkStatusCode(t, resp, 200)

		got, ok := body.(pomodorosResponse)
		if !ok {
			t.Fatal("Type Assertion failed")
		}

		want := pomodorosResponse{
			Pomodoros: []*pomodoroResponse{
				{
					ID: 1,
					Task: &taskResponse{
						ID:                  1,
						Title:               "タスク1",
						ExpectedPomodoroNum: 0,
						ActualPomodoroNum:   2,
						DueOn:               "2021-01-01T00:00:00Z",
						IsCompleted:         false,
						CompletedOn:         "",
						CreatedAt:           "2021-01-01T00:00:00Z",
						UpdatedAt:           "2021-01-01T00:00:00Z",
					},
					CreatedAt: "2021-08-31T01:02:03Z",
				},
			},
		}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("getPomodoros response mismatch (-got +want):\n%s", diff)
		}
	})
}

func TestDeletePomodoro(t *testing.T) {
	setupTestDB(t)
	setupTestPomodoros()
	t.Cleanup(teardownTestDB)
	t.Run("ポモドーロ記録 ID が1の記録を削除する", func(t *testing.T) {
		resp, _ := doTestRequest(t, "DELETE", "/pomodoros/1", nil, nil, "")

		checkStatusCode(t, resp, 204)
	})
}

func TestGetRestCount(t *testing.T) {
	setupTestDB(t)
	t.Cleanup(teardownTestDB)
	t.Run("次の15分休憩までのカウントを取得する", func(t *testing.T) {
		resp, body := doTestRequest(t, "GET", "/pomodoros/rest-count", nil, nil, "restCountResponse")

		checkStatusCode(t, resp, 200)

		got, ok := body.(restCountResponse)
		if !ok {
			t.Fatal("type assertion error")
		}

		want := restCountResponse{
			RestCount: 4,
		}

		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("restCountResponse mismatch (-got +want):\n%s", diff)
		}
	})
}

func BenchmarkPostPomodoros(b *testing.B) {
	setupTestDB(b)
	setupTestPomodoros()
	b.Cleanup(teardownTestDB)
	for i := 0; i < b.N; i++ {
		reqBody := strings.NewReader(`{ "taskID": 1 }`)
		_, _ = doTestRequest(b, "POST", "/pomodoros", nil, reqBody, "pomodoroResponse")
	}
}

func BenchmarkGetPomodoros(b *testing.B) {
	setupTestDB(b)
	setupTestPomodoros()
	b.Cleanup(teardownTestDB)
	for i := 0; i < b.N; i++ {
		_, _ = doTestRequest(b, "GET", "/pomodoros", nil, nil, "pomodorosResponse")
	}
}

func BenchmarkGetRestCount(b *testing.B) {
	setupTestDB(b)
	b.Cleanup(teardownTestDB)

	for i := 0; i < b.N; i++ {
		_, _ = doTestRequest(b, "GET", "/pomodoros/rest-count", nil, nil, "restCountResponse")
	}
}

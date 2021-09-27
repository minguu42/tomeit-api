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

		if resp.StatusCode != 200 {
			t.Error("Status code should be 200, but", resp.StatusCode)
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
				ActualPomodoroNum:   0,
				DueOn:               "2021-01-01T00:00:00Z",
				IsCompleted:         false,
			},
		}

		if diff := cmp.Diff(got, want, pomodoroResponseCmpOpts); diff != "" {
			t.Errorf("pomodoroResponse mismatch (-got +want):\n%s", diff)
		}
	})
}

func setupTestPomodoros() {
	setupTestTasks()
	const createPomodoro1 = `INSERT INTO pomodoros (user_id, task_id, created_at) VALUES (1, 1, '2021-08-31 01:02:03')`
	const createPomodoro2 = `INSERT INTO pomodoros (user_id, task_id, created_at) VALUES (1, 1, '2021-09-01 06:07:08')`

	testDB.Exec(createPomodoro1)
	testDB.Exec(createPomodoro2)
}

func TestGetPomodoros(t *testing.T) {
	setupTestDB(t)
	setupTestPomodoros()
	t.Cleanup(teardownTestDB)
	t.Run("ポモドーロ記録を一覧取得する", func(t *testing.T) {
		resp, body := doTestRequest(t, "GET", "/pomodoros", nil, nil, "pomodorosResponse")

		if resp.StatusCode != 200 {
			t.Error("Status code should be 200, but", resp.StatusCode)
		}

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
						ActualPomodoroNum:   0,
						DueOn:               "2021-01-01T00:00:00Z",
						IsCompleted:         false,
					},
					CreatedAt: "2021-08-31T01:02:03Z",
				},
				{
					ID: 2,
					Task: &taskResponse{
						ID:                  1,
						Title:               "タスク1",
						ExpectedPomodoroNum: 0,
						ActualPomodoroNum:   0,
						DueOn:               "2021-01-01T00:00:00Z",
						IsCompleted:         false,
					},
					CreatedAt: "2021-09-01T06:07:08Z",
				},
			},
		}

		if diff := cmp.Diff(got.Pomodoros[0], want.Pomodoros[0], pomodoroResponseCmpOpts); diff != "" {
			t.Errorf("pomodorosResponse mismatch (-got +want):\n%s", diff)
		}
		if diff := cmp.Diff(got.Pomodoros[1], want.Pomodoros[1], pomodoroResponseCmpOpts); diff != "" {
			t.Errorf("pomodorosResponse mismatch (-got +want):\n%s", diff)
		}
	})
}

//
//func TestGetRestCount(t *testing.T) {
//	t.Run("次の15分休憩までのカウントを取得する", func(t *testing.T) {
//		setupTestDB(t)
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
//		if body.RestCount != 4 {
//			t.Error("restCount should be 4, but", body.RestCount)
//		}
//
//		shutdownTestDB(t)
//	})
//}
//
//func BenchmarkPostPomodoros(b *testing.B) {
//	setupTestDB(b)
//	setupTestPomodoros(b)
//	defer shutdownTestDB(b)
//
//	for i := 0; i < b.N; i++ {
//		reqBody := strings.NewReader(`{ "taskID": 1 }`)
//		req, _ := http.NewRequest("POST", testUrl+"/pomodoros", reqBody)
//
//		_, _ = testClient.Do(req)
//	}
//}
//
//func BenchmarkGetPomodoros(b *testing.B) {
//	setupTestDB(b)
//	setupTestGetPomodoros(b)
//	defer shutdownTestDB(b)
//
//	for i := 0; i < b.N; i++ {
//		req, _ := http.NewRequest("GET", testUrl+"/pomodoros", nil)
//
//		_, _ = testClient.Do(req)
//	}
//}
//
//func BenchmarkGetRestCount(b *testing.B) {
//	setupTestDB(b)
//	defer shutdownTestDB(b)
//
//	for i := 0; i < b.N; i++ {
//		req, _ := http.NewRequest("GET", testUrl+"/pomodoros/rest-count", nil)
//
//		_, _ = testClient.Do(req)
//	}
//}

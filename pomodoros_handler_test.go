package tomeit

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestPostPomodoroLog(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		reqBody := strings.NewReader(`{"taskId": 1}`)
		req, err := http.NewRequest("POST", url+"/pomodoros/logs", reqBody)
		if err != nil {
			t.Errorf("Create request failed: %v", err)
		}
		req.Header.Add("Authorization", "someIdToken")

		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("Do request failed: %v", err)
		}

		bytes, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			t.Errorf("Read response body failed: %v", err)
		}

		log.Println(string(bytes))

		var body pomodoroLogResponse
		if err := json.Unmarshal(bytes, &body); err != nil {
			t.Errorf("Unmarshal json failed: %v", err)
		}

		if resp.StatusCode != 201 {
			t.Errorf("Status code should be 201, but %v", resp.StatusCode)
		}
	})
}

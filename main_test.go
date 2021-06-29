package tomeit

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	_ "github.com/go-sql-driver/mysql"
)

var (
	testClient *http.Client
	testUrl    string
)

func TestMain(m *testing.M) {
	db := OpenDB("mysql", "test:password@tcp(localhost:13306)/db_test?parseTime=true")
	defer CloseDB(db)

	r := chi.NewRouter()

	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello!"))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	testUrl = ts.URL
	testClient = &http.Client{}

	m.Run()
}

func TestExample(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		req, err := http.NewRequest("GET", testUrl+"/", nil)
		if err != nil {
			t.Errorf("Create request failed: %v", err)
		}

		resp, err := testClient.Do(req)
		if err != nil {
			t.Errorf("Do request failed: %v", err)
		}

		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Read response failed: %v", err)
		}
		if err := resp.Body.Close(); err != nil {
			t.Errorf("Close response body failed: %v", err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("Status code should be 201, but %v", resp.StatusCode)
		}
		if string(bytes) != "Hello!" {
			t.Errorf("Response Body should be Hello!, but %v", string(bytes))
		}
	})
}

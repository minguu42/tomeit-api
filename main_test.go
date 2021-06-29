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
			t.Error("Close response body failed:", err)
		}

		if resp.StatusCode != 200 {
			t.Error("Status code should be 201, but", resp.StatusCode)
		}
		if string(bytes) != "Hello!" {
			t.Error("Response Body should be Hello!, but", string(bytes))
		}
	})
}

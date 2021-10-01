package tomeit

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/go-chi/render"
)

type pomodoroResponse struct {
	ID        int           `json:"id"`
	Task      *taskResponse `json:"task"`
	CreatedAt string        `json:"createdAt"`
}

func newPomodoroResponse(p *Pomodoro, db dbInterface) *pomodoroResponse {
	r := pomodoroResponse{
		ID:        p.ID,
		Task:      newTaskResponse(&p.Task, db),
		CreatedAt: p.CreatedAt.Format(time.RFC3339),
	}
	return &r
}

func (p *pomodoroResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type postPomodorosRequest struct {
	TaskID int `json:"taskID"`
}

func (p *postPomodorosRequest) Bind(r *http.Request) error {
	if p.TaskID <= 0 {
		return errors.New("missing required taskID field or taskID is a negative number")
	}
	return nil
}

func postPomodoros(db dbInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody := &postPomodorosRequest{}
		if err := render.Bind(r, reqBody); err != nil {
			log.Println("render.Bind failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		user := r.Context().Value(userKey).(*User)

		pomodoroID, err := db.createPomodoro(user.ID, reqBody.TaskID)
		if err != nil {
			log.Println("db.createPomodoro failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		pomodoro, err := db.getPomodoroByID(pomodoroID)
		if err != nil {
			log.Println("db.getPomodoroByID failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		if err := db.decrementRestCount(user); err != nil {
			log.Println("db.decrementRestCount failed:", err)
			_ = render.Render(w, r, internalServerError(err))
			return
		}

		scheme := "http://"
		if r.TLS != nil {
			scheme = "https://"
		}
		w.Header().Set("Location", scheme+r.Host+r.URL.Path+"/"+strconv.Itoa(pomodoro.ID))
		w.WriteHeader(201)
		if err = render.Render(w, r, newPomodoroResponse(pomodoro, db)); err != nil {
			log.Println("render.Render failed:", err)
			_ = render.Render(w, r, internalServerError(err))
			return
		}
	}
}

func deletePomodoro(db dbInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pomodoroID, err := strconv.ParseInt(chi.URLParam(r, "pomodoroID"), 10, 64)
		if err != nil {
			log.Println("strconv.ParseInt failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		user := r.Context().Value(userKey).(*User)

		pomodoro, err := db.getPomodoroByID(int(pomodoroID))
		if err != nil {
			log.Println("db.getPomodoroByID failed:", err)
			_ = render.Render(w, r, notFoundError(err))
			return
		}
		if !user.hasPomodoro(pomodoro) {
			log.Println("user does not have this pomodoro")
			_ = render.Render(w, r, authorizationError(errors.New("you do not have this pomodoro")))
			return
		}

		if err := db.deletePomodoro(pomodoro); err != nil {
			log.Println("db.deletePomodoro failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		w.WriteHeader(204)
	}
}

type pomodorosResponse struct {
	Pomodoros []*pomodoroResponse `json:"pomodoros"`
}

func newPomodorosResponse(pomodoroRecords []Pomodoro, db dbInterface) *pomodorosResponse {
	var ps []*pomodoroResponse
	for _, p := range pomodoroRecords {
		ps = append(ps, newPomodoroResponse(&p, db))
	}
	return &pomodorosResponse{Pomodoros: ps}
}

func (ps *pomodorosResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func getPomodoros(db dbInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var options getPomodorosOptions
		createdOnStr := r.URL.Query().Get("createdOn")
		createdOn, err := time.Parse(time.RFC3339, createdOnStr)
		if err == nil {
			options.createdOnExists = true
			options.createdOn = createdOn
		}

		user := r.Context().Value(userKey).(*User)

		pomodoros, err := db.getPomodorosByUser(user, &options)
		if err != nil {
			log.Println("db.getPomodorosByUser failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		if err := render.Render(w, r, newPomodorosResponse(pomodoros, db)); err != nil {
			log.Println("render.Render failed:", err)
			_ = render.Render(w, r, internalServerError(err))
			return
		}
	}
}

type restCountResponse struct {
	RestCount int `json:"restCount"`
}

func (c *restCountResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func getRestCount(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userKey).(*User)

	if err := render.Render(w, r, &restCountResponse{RestCount: user.RestCount}); err != nil {
		log.Println("render.Render failed:", err)
		_ = render.Render(w, r, internalServerError(err))
		return
	}
}

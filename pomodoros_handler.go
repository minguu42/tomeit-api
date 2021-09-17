package tomeit

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

type pomodoroResponse struct {
	ID          int64         `json:"id"`
	Task        *taskResponse `json:"task"`
	CompletedAt string        `json:"completedAt"`
	CreatedAt   string        `json:"createdAt"`
}

func newPomodoroResponse(p *pomodoro, db dbInterface) *pomodoroResponse {
	r := pomodoroResponse{
		ID:          p.id,
		Task:        newTaskResponse(p.task, db),
		CompletedAt: p.completedAt.Format(time.RFC3339),
		CreatedAt:   p.createdAt.Format(time.RFC3339),
	}
	return &r
}

func (p *pomodoroResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type pomodorosResponse struct {
	Pomodoros []*pomodoroResponse `json:"pomodoros"`
}

func newPomodorosResponse(pomodoroRecords []*pomodoro, db dbInterface) *pomodorosResponse {
	var ps []*pomodoroResponse
	for _, p := range pomodoroRecords {
		ps = append(ps, newPomodoroResponse(p, db))
	}
	return &pomodorosResponse{Pomodoros: ps}
}

func (ps *pomodorosResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type postPomodorosRequest struct {
	TaskID int64 `json:"taskID"`
}

func (p *postPomodorosRequest) Bind(r *http.Request) error {
	if p.TaskID <= 0 {
		return errors.New("missing required taskID field or taskID is a negative number")
	}
	return nil
}

func PostPomodoros(db dbInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody := &postPomodorosRequest{}
		if err := render.Bind(r, reqBody); err != nil {
			log.Println("render.Bind failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		user := r.Context().Value(userKey).(*user)

		pomodoroID, err := db.createPomodoro(user.id, reqBody.TaskID)
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

		render.Status(r, http.StatusCreated)
		if err = render.Render(w, r, newPomodoroResponse(pomodoro, db)); err != nil {
			log.Println("render.Render failed:", err)
			_ = render.Render(w, r, internalServerError(err))
			return
		}
	}
}

func GetPomodoros(db dbInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		existCompletedOn := true
		completedOnStr := r.URL.Query().Get("completed-on")
		completedOn, err := time.Parse(time.RFC3339, completedOnStr)
		if err != nil {
			if completedOnStr == "" {
				existCompletedOn = false
			} else {
				_ = render.Render(w, r, badRequestError(errors.New("completed-on value is invalid")))
				return
			}
		}

		options := getPomodorosOptions{
			existCompletedOn: existCompletedOn,
			completedOn:      completedOn,
		}

		user := r.Context().Value(userKey).(*user)

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

func GetRestCount(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userKey).(*user)

	if err := render.Render(w, r, &restCountResponse{RestCount: user.restCount}); err != nil {
		log.Println("render.Render failed:", err)
		_ = render.Render(w, r, internalServerError(err))
		return
	}
}

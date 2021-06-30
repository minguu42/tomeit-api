package tomeit

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

type pomodoroLogRequest struct {
	TaskID int64 `json:"taskID"`
}

func (p *pomodoroLogRequest) Bind(r *http.Request) error {
	if p.TaskID == 0 {
		return errors.New("missing required taskID field")
	}
	return nil
}

type pomodoroLogResponse struct {
	ID        int64         `json:"id"`
	Task      *taskResponse `json:"task"`
	CreatedAt string        `json:"createdAt"`
}

func newPomodoroLogResponse(p *pomodoroLog) *pomodoroLogResponse {
	r := pomodoroLogResponse{
		ID:        p.id,
		Task:      newTaskResponse(p.task),
		CreatedAt: p.createdAt.Format(time.RFC3339),
	}
	return &r
}

func (p *pomodoroLogResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func PostPomodoroLog(db dbInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &pomodoroLogRequest{}
		if err := render.Bind(r, data); err != nil {
			log.Println("bind failed:", err)
			_ = render.Render(w, r, errInvalidRequest(err))
			return
		}

		user := r.Context().Value("user").(*user)

		pomodoroLogID, err := db.createPomodoroLog(user.id, data.TaskID)
		if err != nil {
			log.Println("createPomodoroLog failed:", err)
			_ = render.Render(w, r, errInvalidRequest(err))
			return
		}

		pomodoroLog, err := db.getPomodoroLogByID(pomodoroLogID)
		if err != nil {
			log.Println("getPomodoroLogByID failed:", err)
			_ = render.Render(w, r, errInvalidRequest(err))
			return
		}

		render.Status(r, http.StatusCreated)
		if err = render.Render(w, r, newPomodoroLogResponse(pomodoroLog)); err != nil {
			log.Println("render failed:", err)
			_ = render.Render(w, r, errRender(err))
			return
		}
	}
}

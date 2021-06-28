package tomeit

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

type pomodoroLogRequest struct {
	TaskId int64 `json:"taskId"`
}

func (p *pomodoroLogRequest) Bind(r *http.Request) error {
	if p.TaskId <= 0 {
		return errors.New("taskId should be positive num")
	}
	return nil
}

type pomodoroLogResponse struct {
	Id        int64         `json:"id"`
	Task      *taskResponse `json:"task"`
	CreatedAt string        `json:"createdAt"`
}

func newPomodoroLogResponse(pomodoroLog *pomodoroLog) *pomodoroLogResponse {
	resp := pomodoroLogResponse{
		Id:        pomodoroLog.id,
		Task:      newTaskResponse(pomodoroLog.task),
		CreatedAt: pomodoroLog.createdAt.Format(time.RFC3339),
	}

	return &resp
}

func (resp pomodoroLogResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func PostPomodoroLog(w http.ResponseWriter, r *http.Request) {
	data := &pomodoroLogRequest{}
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, errInvalidRequest(err))
		return
	}

	user := r.Context().Value("user").(User)

	id, err := createPomodoroLog(user.id, data.TaskId)
	if err != nil {
		_ = render.Render(w, r, errInvalidRequest(err))
		return
	}

	pomodoroLog, err := getPomodoroLogById(id)
	if err != nil {
		_ = render.Render(w, r, errInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusCreated)
	if err := render.Render(w, r, newPomodoroLogResponse(pomodoroLog)); err != nil {
		_ = render.Render(w, r, errRender(err))
	}
}

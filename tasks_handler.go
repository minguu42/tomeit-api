package tomeit

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

type taskRequest struct {
	*Task
}

func (t *taskRequest) Bind(r *http.Request) error {
	if t.Task == nil {
		return errors.New("missing required Task fields")
	}
	if t.UserID == 0 {
		return errors.New("missing required UserID fields")
	}
	if t.Name == "" {
		return errors.New("missing required Name fields")
	}
	return nil
}

type taskResponse struct {
	*Task
}

func newTaskResponse(task *Task) *taskResponse {
	resp := &taskResponse{task}
	return resp
}

func (t taskResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func PostTask(w http.ResponseWriter, r *http.Request) {
	data := &taskRequest{}
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, errInvalidRequest(err))
		return
	}

	task := data.Task

	createdTaskId, err := createTask(task.UserID, task.Name, task.Priority, task.Deadline)
	if err != nil {
		_ = render.Render(w, r, errInvalidRequest(err))
		return
	}

	createdTask, err := getTaskById(createdTaskId)
	if err != nil {
		_ = render.Render(w, r, errInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusCreated)
	_ = render.Render(w, r, newTaskResponse(&createdTask))
}

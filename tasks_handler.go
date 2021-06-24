package tomeit

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

type taskRequest struct {
	Name     string `json:"name"`
	Priority int    `json:"priority,omitempty"`
	Deadline string `json:"deadline,omitempty"`
}

func (t *taskRequest) Bind(r *http.Request) error {
	if t.Name == "" {
		return errors.New("missing required Name fields")
	}
	if t.Deadline == "" {
		t.Deadline = "0001-01-01"
	}
	return nil
}

type taskResponse struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Priority      int    `json:"priority"`
	Deadline      string `json:"deadline"`
	IsDone        bool   `json:"isDone"`
	PomodoroCount int    `json:"pomodoroCount"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

func newTaskResponse(task *Task) *taskResponse {
	resp := taskResponse{
		Id:            task.id,
		Name:          task.name,
		Priority:      task.priority,
		Deadline:      task.deadline.Format("2006-01-02"),
		IsDone:        task.isDone,
		PomodoroCount: 0,
		CreatedAt:     task.createdAt.Format(time.RFC3339),
		UpdatedAt:     task.updatedAt.Format(time.RFC3339),
	}
	return &resp
}

func (t taskResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func PostTask(w http.ResponseWriter, r *http.Request) {
	data := &taskRequest{}
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, invalidRequestErr(err))
		return
	}

	deadline, err := time.Parse("2006-01-02", data.Deadline)
	if err != nil {
		_ = render.Render(w, r, invalidRequestErr(err))
		return
	}
	createdTaskId, err := createTask(1, data.Name, data.Priority, deadline)
	if err != nil {
		_ = render.Render(w, r, invalidRequestErr(err))
		return
	}

	createdTask, err := getTaskById(createdTaskId)
	if err != nil {
		_ = render.Render(w, r, invalidRequestErr(err))
		return
	}

	render.Status(r, http.StatusCreated)
	_ = render.Render(w, r, newTaskResponse(&createdTask))
}

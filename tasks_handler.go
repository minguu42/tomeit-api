package tomeit

import (
	"errors"
	"log"
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
		return errors.New("missing required name fields")
	}
	if t.Deadline == "" {
		t.Deadline = "0001-01-01"
	}
	return nil
}

type taskResponse struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Priority      int    `json:"priority"`
	Deadline      string `json:"deadline"`
	IsDone        bool   `json:"isDone"`
	PomodoroCount int    `json:"pomodoroCount"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

func (t *taskResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func newTaskResponse(t *task) *taskResponse {
	r := taskResponse{
		ID:            t.id,
		Name:          t.name,
		Priority:      t.priority,
		Deadline:      t.deadline.Format("2006-01-02"),
		IsDone:        t.isDone,
		PomodoroCount: 0,
		CreatedAt:     t.createdAt.Format(time.RFC3339),
		UpdatedAt:     t.updatedAt.Format(time.RFC3339),
	}
	return &r
}

func PostTask(db dbInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &taskRequest{}
		if err := render.Bind(r, data); err != nil {
			log.Println("bind failed:", err)
			_ = render.Render(w, r, errInvalidRequest(err))
			return
		}

		deadline, err := time.Parse("2006-01-02", data.Deadline)
		if err != nil {
			log.Println("parse failed:", err)
			_ = render.Render(w, r, errInvalidRequest(err))
			return
		}

		user := r.Context().Value("user").(*user)

		taskID, err := db.createTask(user.id, data.Name, data.Priority, deadline)
		if err != nil {
			log.Println("createTask failed:", err)
			_ = render.Render(w, r, errInvalidRequest(err))
			return
		}

		task, err := db.getTaskByID(taskID)
		if err != nil {
			log.Println("getTaskByID failed:", err)
			_ = render.Render(w, r, errInvalidRequest(err))
			return
		}

		render.Status(r, http.StatusCreated)
		if err = render.Render(w, r, newTaskResponse(task)); err != nil {
			log.Println("render failed:", err)
			_ = render.Render(w, r, errRender(err))
			return
		}
	}
}

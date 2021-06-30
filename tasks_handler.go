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

func (t *taskResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type tasksResponse struct {
	Tasks []*taskResponse `json:"tasks"`
}

func newTasksResponse(tasks []*task) *tasksResponse {
	var ts []*taskResponse
	for _, t := range tasks {
		ts = append(ts, newTaskResponse(t))
	}
	return &tasksResponse{Tasks: ts}
}

func (ts *tasksResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
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

func GetTasks(db dbInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*user)

		tasks, err := db.getTasksByUser(user)
		if err != nil {
			log.Println("getTasksByUser failed:", err)
			_ = render.Render(w, r, errNotFound())
			return
		}

		if err := render.Render(w, r, newTasksResponse(tasks)); err != nil {
			log.Println("render failed:", err)
			_ = render.Render(w, r, errRender(err))
			return
		}
	}
}

func GetTasksDone(db dbInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*user)

		tasks, err := db.getDoneTasksByUser(user)
		if err != nil {
			log.Println("getTasksByUser failed:", err)
			_ = render.Render(w, r, errNotFound())
			return
		}

		if err := render.Render(w, r, newTasksResponse(tasks)); err != nil {
			log.Println("render failed:", err)
			_ = render.Render(w, r, errRender(err))
			return
		}
	}
}

func PutTaskDone(db dbInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskIDStr := chi.URLParam(r, "taskID")
		if taskIDStr == "" {
			_ = render.Render(w, r, errNotFound())
			return
		}

		taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
		if err != nil {
			log.Println("parseInt failed:", err)
			_ = render.Render(w, r, errInvalidRequest(err))
			return
		}

		user := r.Context().Value("user").(*user)

		if !hasUserTask(db, taskID, user) {
			_ = render.Render(w, r, errAuthenticate(errors.New("you do not have this task")))
			return
		}

		if err := db.doneTask(taskID); err != nil {
			log.Println("doneTask failed:", err)
			_ = render.Render(w, r, errInvalidRequest(err))
			return
		}

		task, err := db.getTaskByID(taskID)
		if err != nil {
			log.Println("getTaskByID failed:", err)
			_ = render.Render(w, r, errInvalidRequest(err))
			return
		}

		if err = render.Render(w, r, newTaskResponse(task)); err != nil {
			log.Println("render failed:", err)
			_ = render.Render(w, r, errRender(err))
			return
		}
	}
}

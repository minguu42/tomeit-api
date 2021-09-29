package tomeit

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/go-chi/render"
)

type taskResponse struct {
	ID                  int    `json:"id"`
	Title               string `json:"title"`
	ExpectedPomodoroNum int    `json:"expectedPomodoroNum"`
	ActualPomodoroNum   int    `json:"actualPomodoroNum"`
	DueOn               string `json:"dueOn"`
	IsCompleted         bool   `json:"isCompleted"`
	CompletedOn         string `json:"completedOn"`
	CreatedAt           string `json:"createdAt"`
	UpdatedAt           string `json:"updatedAt"`
}

func newTaskResponse(t *Task, db dbInterface) *taskResponse {
	c, err := db.getActualPomodoroNumByID(t.ID)
	if err != nil {
		c = 0
	}

	var dueOn string
	if t.DueOn != nil {
		dueOn = t.DueOn.Format(time.RFC3339)
	}

	var completedOn string
	if t.CompletedOn != nil {
		completedOn = t.CompletedOn.Format(time.RFC3339)
	}

	r := taskResponse{
		ID:                  t.ID,
		Title:               t.Title,
		ExpectedPomodoroNum: t.ExpectedPomodoroNum,
		ActualPomodoroNum:   c,
		DueOn:               dueOn,
		IsCompleted:         t.IsCompleted,
		CompletedOn:         completedOn,
		CreatedAt:           t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           t.UpdatedAt.Format(time.RFC3339),
	}
	return &r
}

func (t *taskResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type postTasksRequest struct {
	Title               string `json:"title"`
	ExpectedPomodoroNum int    `json:"expectedPomodoroNum,omitempty"`
	DueOn               string `json:"dueOn,omitempty"`
}

func (p *postTasksRequest) Bind(r *http.Request) error {
	if p.Title == "" {
		return errors.New("missing required title field")
	}
	if p.DueOn == "" {
		p.DueOn = "0001-01-01T00:00:00Z"
	}
	return nil
}

func postTasks(db dbInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody := &postTasksRequest{}
		if err := render.Bind(r, reqBody); err != nil {
			log.Println("render.Bind failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		dueAt, err := time.Parse(time.RFC3339, reqBody.DueOn)
		if err != nil {
			log.Println("time.Parse failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		user := r.Context().Value(userKey).(*User)

		task, err := db.createTask(user.ID, reqBody.Title, reqBody.ExpectedPomodoroNum, dueAt)
		if err != nil {
			log.Println("db.createTask failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		scheme := "http://"
		if r.TLS != nil {
			scheme = "https://"
		}
		w.Header().Set("Location", scheme+r.Host+r.URL.Path+"/"+strconv.Itoa(task.ID))
		w.WriteHeader(201)
		if err = render.Render(w, r, newTaskResponse(task, db)); err != nil {
			log.Println("render.Render failed:", err)
			_ = render.Render(w, r, internalServerError(err))
			return
		}
	}
}

type tasksResponse struct {
	Tasks []taskResponse `json:"tasks"`
}

func newTasksResponse(tasks []Task, db dbInterface) *tasksResponse {
	var ts []taskResponse
	for _, t := range tasks {
		ts = append(ts, *newTaskResponse(&t, db))
	}
	return &tasksResponse{Tasks: ts}
}

func (ts *tasksResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func getTasks(db dbInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var options getTasksOptions

		isCompletedStr := r.URL.Query().Get("isCompleted")
		if isCompletedStr == "true" {
			options.isCompletedExists = true
			options.isCompleted = true
		} else if isCompletedStr == "false" {
			options.isCompletedExists = false
			options.isCompleted = false
		}

		completedOnStr := r.URL.Query().Get("completedOn")
		if completedOn, err := time.Parse(time.RFC3339, completedOnStr); err == nil {
			options.completedOnExists = true
			options.completedOn = completedOn
		}

		user := r.Context().Value(userKey).(*User)

		tasks, err := db.getTasksByUser(user, &options)
		if err != nil {
			log.Println("db.getTasksByUser failed:", err)
			_ = render.Render(w, r, internalServerError(err))
			return
		}

		if err := render.Render(w, r, newTasksResponse(tasks, db)); err != nil {
			log.Println("render.Render failed:", err)
			_ = render.Render(w, r, internalServerError(err))
			return
		}
	}
}

type patchTaskRequest struct {
	IsCompleted string `json:"isCompleted"`
}

func (p *patchTaskRequest) Bind(r *http.Request) error {
	if p.IsCompleted != "true" && p.IsCompleted != "false" {
		return fmt.Errorf("isCompleted value is wrong")
	}
	return nil
}

func patchTask(db dbInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskID, err := strconv.ParseInt(chi.URLParam(r, "taskID"), 10, 64)
		if err != nil {
			log.Println("strconv.ParseInt failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		user := r.Context().Value(userKey).(*User)

		task, err := db.getTaskByID(int(taskID))
		if err != nil {
			log.Println("db.getTaskByID failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		if !user.hasTask(task) {
			log.Println("user does not have a task")
			_ = render.Render(w, r, AuthorizationError(errors.New("task's userID does not match your userID")))
			return
		}

		data := &patchTaskRequest{}
		if err := render.Bind(r, data); err != nil {
			log.Println("render.Bind failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		task.IsCompleted, _ = strconv.ParseBool(data.IsCompleted)

		if err := db.updateTask(task); err != nil {
			log.Println("db.updateTask failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		task, err = db.getTaskByID(task.ID)
		if err != nil {
			log.Println("db.getTaskByID failed:", err)
			_ = render.Render(w, r, internalServerError(err))
			return
		}

		if err := render.Render(w, r, newTaskResponse(task, db)); err != nil {
			log.Println("render.Render failed:", err)
			_ = render.Render(w, r, internalServerError(err))
			return
		}
	}
}

func deleteTask(db dbInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskID, err := strconv.ParseInt(chi.URLParam(r, "taskID"), 10, 64)
		if err != nil {
			log.Println("strconv.ParseInt failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		user := r.Context().Value(userKey).(*User)

		task, err := db.getTaskByID(int(taskID))
		if err != nil {
			log.Println("db.getTaskByID failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		if !user.hasTask(task) {
			log.Println("user does not have a task")
			_ = render.Render(w, r, AuthorizationError(errors.New("task's userID does not match your userID")))
			return
		}

		if err := db.deleteTask(task); err != nil {
			log.Println("db.deleteTask failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		w.WriteHeader(204)
	}
}

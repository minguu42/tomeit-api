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

type taskResponse struct {
	ID                     int64  `json:"id"`
	Title                  string `json:"title"`
	ExpectedPomodoroNumber int    `json:"expectedPomodoroNumber"`
	ActualPomodoroNumber   int    `json:"actualPomodoroNumber"`
	DueOn                  string `json:"dueOn"`
	IsCompleted            bool   `json:"isCompleted"`
	CompletedAt            string `json:"completedOn"`
	CreatedAt              string `json:"createdAt"`
	UpdatedAt              string `json:"updatedAt"`
}

func newTaskResponse(t *task, db dbInterface) *taskResponse {
	c, err := db.getActualPomodoroNumberByID(t.id)
	if err != nil {
		c = 0
	}

	r := taskResponse{
		ID:                     t.id,
		Title:                  t.title,
		ExpectedPomodoroNumber: t.expectedPomodoroNumber,
		ActualPomodoroNumber:   c,
		DueOn:                  t.dueOn.Format(time.RFC3339),
		IsCompleted:            t.isCompleted,
		CompletedAt:            t.completedAt.Format(time.RFC3339),
		CreatedAt:              t.createdAt.Format(time.RFC3339),
		UpdatedAt:              t.updatedAt.Format(time.RFC3339),
	}
	return &r
}

func (t *taskResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type tasksResponse struct {
	Tasks []*taskResponse `json:"tasks"`
}

func newTasksResponse(tasks []*task, db dbInterface) *tasksResponse {
	var ts []*taskResponse
	for _, t := range tasks {
		ts = append(ts, newTaskResponse(t, db))
	}
	return &tasksResponse{Tasks: ts}
}

func (ts *tasksResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type postTaskRequest struct {
	Title                  string `json:"title"`
	ExpectedPomodoroNumber int    `json:"expectedPomodoroNumber,omitempty"`
	DueOn                  string `json:"dueOn,omitempty"`
}

func (p *postTaskRequest) Bind(r *http.Request) error {
	if p.Title == "" {
		return errors.New("missing required title field")
	}
	if p.DueOn == "" {
		p.DueOn = "0001-01-01T00:00:00Z"
	}
	return nil
}

func PostTask(db dbInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody := &postTaskRequest{}
		if err := render.Bind(r, reqBody); err != nil {
			log.Println("render.Bind failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		dueOn, err := time.Parse(time.RFC3339, reqBody.DueOn)
		if err != nil {
			log.Println("time.Parse failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		user := r.Context().Value(userKey).(*user)

		taskID, err := db.createTask(user.id, reqBody.Title, reqBody.ExpectedPomodoroNumber, dueOn)
		if err != nil {
			log.Println("db.createTask failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		task, err := db.getTaskByID(taskID)
		if err != nil {
			log.Println("db.getTaskByID failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		render.Status(r, http.StatusCreated)
		if err = render.Render(w, r, newTaskResponse(task, db)); err != nil {
			log.Println("render.Render failed:", err)
			_ = render.Render(w, r, internalServerError(err))
			return
		}
	}
}

func GetTasks(db dbInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		existIsCompleted := true
		var isCompleted bool
		isCompletedStr := r.URL.Query().Get("is-completed")
		if isCompletedStr == "" {
			existIsCompleted = false
		} else if isCompletedStr == "true" {
			isCompleted = true
		} else if isCompletedStr == "false" {
			isCompleted = false
		} else {
			_ = render.Render(w, r, badRequestError(errors.New("is-completed value is invalid")))
			return
		}

		existCompletedAt := true
		completedAtStr := r.URL.Query().Get("completed-at")
		completedAt, err := time.Parse(time.RFC3339, completedAtStr)
		if err != nil {
			if completedAtStr == "" {
				existCompletedAt = false
			} else {
				_ = render.Render(w, r, badRequestError(errors.New("completed-at value is invalid")))
				return
			}
		}

		options := getTasksOptions{
			existIsCompleted: existIsCompleted,
			isCompleted:      isCompleted,
			existCompletedAt: existCompletedAt,
			completedAt:      completedAt,
		}

		user := r.Context().Value(userKey).(*user)

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
	IsCompleted bool `json:"isCompleted"`
}

func (p *patchTaskRequest) Bind(r *http.Request) error {
	return nil
}

func PatchTask(db dbInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskID, err := strconv.ParseInt(chi.URLParam(r, "task-id"), 10, 64)
		if err != nil {
			log.Println("strconv.ParseInt failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		user := r.Context().Value(userKey).(*user)

		task, err := db.getTaskByID(taskID)
		if err != nil {
			log.Println("db.getTaskByID failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}
		if user.id != task.user.id {
			log.Println("user.id != task.user.id")
			_ = render.Render(w, r, AuthorizationError(errors.New("task's userID does not match your userID")))
			return
		}

		data := &patchTaskRequest{}
		if err := render.Bind(r, data); err != nil {
			log.Println("render.Bind failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		options := updateTaskOptions{
			existIsCompleted: true,
		}

		task.isCompleted = data.IsCompleted
		if err := db.updateTask(task, &options); err != nil {
			log.Println("db.updateTask failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		task, err = db.getTaskByID(task.id)
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

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
	Title                  string `json:"title"`
	ExpectedPomodoroNumber int    `json:"expectedPomodoroNumber,omitempty"`
	DueOn                  string `json:"dueOn,omitempty"`
}

func (t *taskRequest) Bind(r *http.Request) error {
	if t.Title == "" {
		return errors.New("missing required title field")
	}
	if t.DueOn == "" {
		t.DueOn = "0001-01-01T00:00:00Z"
	}
	return nil
}

type taskResponse struct {
	ID                     int64  `json:"id"`
	Title                  string `json:"title"`
	ExpectedPomodoroNumber int    `json:"expectedPomodoroNumber"`
	ActualPomodoroNumber   int    `json:"actualPomodoroNumber"`
	DueOn                  string `json:"dueOn"`
	IsCompleted            bool   `json:"isCompleted"`
	CompletedAt            string `json:"completedAt"`
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

func PostTask(db dbInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &taskRequest{}
		if err := render.Bind(r, data); err != nil {
			log.Println("render.Bind failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		dueOn, err := time.Parse(time.RFC3339, data.DueOn)
		if err != nil {
			log.Println("time.Parse failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		user := r.Context().Value(userKey).(*user)

		taskID, err := db.createTask(user.id, data.Title, data.ExpectedPomodoroNumber, dueOn)
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
		//var isCompleted bool
		isCompletedStr := chi.URLParam(r, "is-completed")
		//if isCompletedStr == "true" {
		//	isCompleted = true
		//} else if isCompletedStr == "false" {
		//	isCompleted = false
		//}
		//
		completedAtStr := chi.URLParam(r, "completed-at")
		//completedAt, err := time.Parse(time.RFC3339, completedAtStr)
		//if err != nil {
		//	log.Println("time.Parse failed:", err)
		//	_ = render.Render(w, r, badRequestError(err))
		//	return
		//}

		user := r.Context().Value(userKey).(*user)

		var tasks []*task
		var err error
		if isCompletedStr == "" && completedAtStr == "" {
			tasks, err = db.getTasksByUser(user)
			if err != nil {
				log.Println("db.getTasksByUser failed:", err)
				_ = render.Render(w, r, errNotFound())
				return
			}
		}

		if err := render.Render(w, r, newTasksResponse(tasks, db)); err != nil {
			log.Println("render.Render failed:", err)
			_ = render.Render(w, r, internalServerError(err))
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
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		user := r.Context().Value(userKey).(*user)

		if !hasUserTask(db, taskID, user) {
			_ = render.Render(w, r, AuthenticationError(errors.New("you do not have this task")))
			return
		}

		if err := db.doneTask(taskID); err != nil {
			log.Println("doneTask failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		task, err := db.getTaskByID(taskID)
		if err != nil {
			log.Println("getTaskByID failed:", err)
			_ = render.Render(w, r, badRequestError(err))
			return
		}

		if err = render.Render(w, r, newTaskResponse(task, db)); err != nil {
			log.Println("render failed:", err)
			_ = render.Render(w, r, renderError(err))
			return
		}
	}
}

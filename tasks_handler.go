package tomeit

import (
	"errors"
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
		return errors.New("missing required Name fields")
	}
	if t.Deadline == "" {
		t.Deadline = "0001-01-01"
	}
	return nil
}

type taskResponse struct {
	Id            int64  `json:"id"`
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

type tasksResponse struct {
	Tasks []*taskResponse `json:"tasks"`
}

func newTasksResponse(tasks []*Task) *tasksResponse {
	var ts []*taskResponse
	for _, t := range tasks {
		ts = append(ts, newTaskResponse(t))
	}
	var resp tasksResponse
	resp.Tasks = ts
	return &resp
}

func (ts tasksResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func PostTask(w http.ResponseWriter, r *http.Request) {
	data := &taskRequest{}
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, invalidRequestErr(err))
		return
	}
	user := r.Context().Value("user").(User)

	deadline, err := time.Parse("2006-01-02", data.Deadline)
	if err != nil {
		_ = render.Render(w, r, invalidRequestErr(err))
		return
	}

	createdTaskId, err := createTask(user.id, data.Name, data.Priority, deadline)
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

func GetUndoneTasks(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(User)

	tasks, err := getUndoneTasksByUserID(user.id)
	if err != nil {
		_ = render.Render(w, r, invalidRequestErr(err))
		return
	}

	_ = render.Render(w, r, newTasksResponse(tasks))
}

func GetDoneTasks(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(User)

	tasks, err := getDoneTasksByUserID(user.id)
	if err != nil {
		_ = render.Render(w, r, invalidRequestErr(err))
		return
	}

	_ = render.Render(w, r, newTasksResponse(tasks))
}

func PutTaskDone(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(User)

	strTaskId := chi.URLParam(r, "taskId")
	if strTaskId == "" {
		_ = render.Render(w, r, invalidRequestErr(errors.New("URL path does not have taskId")))
		return
	}

	taskId, err := strconv.ParseInt(strTaskId, 10, 64)
	if err != nil {
		_ = render.Render(w, r, invalidRequestErr(errors.New("taskId must be number")))
		return
	}

	task, err := getTaskById(taskId)
	if err != nil {
		_ = render.Render(w, r, invalidRequestErr(errors.New("task does not exits")))
		return
	}

	if user.id != task.userId {
		_ = render.Render(w, r, authenticateErr(errors.New("you do not own this task")))
		return
	}

	if err := completeTask(task.id); err != nil {
		_ = render.Render(w, r, unexpectedErr(err))
		return
	}
	task.isDone = true

	_ = render.Render(w, r, newTaskResponse(&task))
}

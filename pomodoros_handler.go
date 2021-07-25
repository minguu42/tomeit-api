package tomeit

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

type pomodoroRecordRequest struct {
	TaskID int64 `json:"taskID"`
}

func (p *pomodoroRecordRequest) Bind(r *http.Request) error {
	if p.TaskID == 0 {
		return errors.New("missing required taskID field")
	}
	return nil
}

type pomodoroRecordResponse struct {
	ID        int64         `json:"id"`
	Task      *taskResponse `json:"task"`
	CreatedAt string        `json:"createdAt"`
}

func newPomodoroRecordResponse(p *pomodoroRecord, db dbInterface) *pomodoroRecordResponse {
	r := pomodoroRecordResponse{
		ID:        p.id,
		Task:      newTaskResponse(p.task, db),
		CreatedAt: p.createdAt.Format(time.RFC3339),
	}
	return &r
}

func (p *pomodoroRecordResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type pomodoroRecordsResponse struct {
	PomodoroRecords []*pomodoroRecordResponse `json:"pomodoroRecords"`
}

func newPomodoroRecordsResponse(pomodoroLogs []*pomodoroRecord, db dbInterface) *pomodoroRecordsResponse {
	var ps []*pomodoroRecordResponse
	for _, p := range pomodoroLogs {
		ps = append(ps, newPomodoroRecordResponse(p, db))
	}
	return &pomodoroRecordsResponse{PomodoroRecords: ps}
}

func (ps *pomodoroRecordsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func PostPomodoroRecord(db dbInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &pomodoroRecordRequest{}
		if err := render.Bind(r, data); err != nil {
			log.Println("bind failed:", err)
			_ = render.Render(w, r, errBadRequest(err))
			return
		}

		user := r.Context().Value(userKey).(*user)

		pomodoroLogID, err := db.createPomodoroRecord(user.id, data.TaskID)
		if err != nil {
			log.Println("createPomodoroRecord failed:", err)
			_ = render.Render(w, r, errBadRequest(err))
			return
		}

		pomodoroLog, err := db.getPomodoroRecordByID(pomodoroLogID)
		if err != nil {
			log.Println("getPomodoroRecordByID failed:", err)
			_ = render.Render(w, r, errBadRequest(err))
			return
		}

		if err := db.decrementRestCount(user); err != nil {
			log.Println("decrementRestCount failed:", err)
			_ = render.Render(w, r, errUnexpectedEvent(err))
			return
		}

		render.Status(r, http.StatusCreated)
		if err = render.Render(w, r, newPomodoroRecordResponse(pomodoroLog, db)); err != nil {
			log.Println("render failed:", err)
			_ = render.Render(w, r, errRender(err))
			return
		}
	}
}

func GetPomodoroRecords(db dbInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(userKey).(*user)

		pomodoroLogs, err := db.getPomodoroRecordsByUser(user)
		if err != nil {
			log.Println("getPomodoroRecordsByUser failed:", err)
			_ = render.Render(w, r, errNotFound())
			return
		}

		if err := render.Render(w, r, newPomodoroRecordsResponse(pomodoroLogs, db)); err != nil {
			log.Println("render failed:", err)
			_ = render.Render(w, r, errRender(err))
			return
		}
	}
}

type todayPomodoroCountResponse struct {
	TodayPomodoroCount int `json:"todayPomodoroCount"`
}

func (resp *todayPomodoroCountResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func GetTodayPomodoroCount(db dbInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(userKey).(*user)

		count, err := db.getTodayPomodoroCount(user)
		if err != nil {
			log.Println("getTodayPomodoroCount failed:", err)
			_ = render.Render(w, r, errBadRequest(err))
			return
		}

		if err := render.Render(w, r, &todayPomodoroCountResponse{TodayPomodoroCount: count}); err != nil {
			log.Println("render failed:", err)
			_ = render.Render(w, r, errRender(err))
			return
		}
	}
}

type restCountResponse struct {
	RestCount int `json:"restCount"`
}

func (c *restCountResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func GetRestCount(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userKey).(*user)

	if err := render.Render(w, r, &restCountResponse{RestCount: user.restCount}); err != nil {
		log.Println("render failed:", err)
		_ = render.Render(w, r, errRender(err))
		return
	}
}

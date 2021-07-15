package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/minguu42/tomeit-api"

	"github.com/go-chi/cors"

	"github.com/go-chi/render"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	firebaseApp := tomeit.InitFirebaseApp()

	db := tomeit.OpenDB("mysql", os.Getenv("DSN"))
	defer tomeit.CloseDB(db)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   strings.Split(os.Getenv("ALLOW_ORIGINS"), ","),
		AllowedMethods:   []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}))
	r.Use(tomeit.UserCtx(db, firebaseApp))

	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", tomeit.PostTask(db))
		r.Get("/", tomeit.GetTasks(db))

		r.Route("/done", func(r chi.Router) {
			r.Get("/", tomeit.GetTasksDone(db))
			r.Put("/{taskID}", tomeit.PutTaskDone(db))
		})
	})
	r.Route("/pomodoros", func(r chi.Router) {
		r.Route("/logs", func(r chi.Router) {
			r.Post("/", tomeit.PostPomodoroLog(db))
			r.Get("/", tomeit.GetPomodoroLogs(db))
		})
		r.Get("/rest/count", tomeit.GetRestCount)
	})

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalln("$PORT must be set")
	}

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalln("ListenAndServe failed:", err)
	}
}

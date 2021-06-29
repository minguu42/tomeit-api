package main

import (
	"log"
	"net/http"
	"os"

	"github.com/minguu42/tomeit-api"

	"github.com/go-chi/cors"

	"github.com/go-chi/render"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	firebaseApp := tomeit.InitFirebaseApp()

	db := tomeit.OpenDB("mysql", os.Getenv("DATABASE_URL"))
	defer tomeit.CloseDB(db)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://tomeit.vercel.app"}, //TODO: 開発後は http://localhost:3000 を除外する.
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
		})
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalln("ListenAndServe failed:", err)
	}
}

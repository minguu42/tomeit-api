package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/cors"

	"github.com/go-chi/render"
	"github.com/minguu42/tomeit-api"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://tomeit.vercel.app"}, //TODO: 開発後は http://localhost:3000 を除外する.
		AllowedMethods:   []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}))

	tomeit.OpenDb()
	defer tomeit.CloseDb()

	tomeit.InitFirebaseApp()

	r.Use(tomeit.UserCtx)
	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", tomeit.PostTask)
		r.Get("/", tomeit.GetUndoneTasks)

		r.Route("/done", func(r chi.Router) {
			r.Get("/", tomeit.GetDoneTasks)
			r.Put("/{taskId}", tomeit.PutTaskDone)
		})
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("ListenAndServe error:", err)
	}
}

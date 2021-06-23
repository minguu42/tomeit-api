package main

import (
	"fmt"
	"net/http"

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

	tomeit.OpenDb()
	defer tomeit.CloseDb()

	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", tomeit.PostTask)
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("ListenAndServe error:", err)
	}
}

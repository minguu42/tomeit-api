package tomeit

import "github.com/go-chi/chi/v5"

func Route(r chi.Router, db dbInterface) {
	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", postTasks(db))
		r.Get("/", getTasks(db))
		r.Patch("/{taskID}", patchTask(db))
		r.Delete("/{taskID}", deleteTask(db))
	})

	r.Route("/pomodoros", func(r chi.Router) {
		//r.Post("/", tomeit.PostPomodoros(db))
		//r.Get("/", tomeit.GetPomodoros(db))

		//r.Get("/rest-count", tomeit.GetRestCount)
	})
}

package routers

import (
	"net/http"

	wrk "github.com/cutlery47/map-reduce/worker/internal/domain/worker"
	"github.com/cutlery47/map-reduce/worker/internal/routers/http/v1/worker"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(wrk wrk.Worker) *chi.Mux {
	var (
		r = chi.NewRouter()
	)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Logger)

		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte("pong"))
		})
		r.Mount("/worker", worker.New(wrk))
	})

	return r
}

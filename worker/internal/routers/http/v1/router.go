package handlers

import (
	"net/http"

	"github.com/cutlery47/map-reduce/worker/internal/domain"
	"github.com/cutlery47/map-reduce/worker/internal/routers/http/v1/mapreduce"
	"github.com/go-chi/chi/v5"
)

func New(wrk domain.Worker) *chi.Mux {
	var (
		r = chi.NewRouter()
	)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("pong"))
	})

	r.Mount("/worker", mapreduce.New(wrk))

	return r
}

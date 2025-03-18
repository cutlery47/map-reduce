package routers

import (
	"net/http"

	"github.com/cutlery47/map-reduce/worker/internal/domain/worker"
	"github.com/cutlery47/map-reduce/worker/internal/routers/http/v1/mapreduce"
	"github.com/go-chi/chi/v5"
)

func New(wrk worker.Worker) *chi.Mux {
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

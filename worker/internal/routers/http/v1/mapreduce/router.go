package mapreduce

import (
	"github.com/cutlery47/map-reduce/worker/internal/domain"
	"github.com/go-chi/chi/v5"
)

func New(wrk domain.Worker) *chi.Mux {
	var (
		r  = chi.NewRouter()
		wr = &workerRoutes{
			wrk: wrk,
		}
	)

	r.Post("/map", wr.handleMap)
	r.Post("/reduce", wr.handleReduce)

	return r
}

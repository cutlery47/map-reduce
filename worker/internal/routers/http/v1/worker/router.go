package worker

import (
	"github.com/cutlery47/map-reduce/worker/internal/domain/worker"
	"github.com/go-chi/chi/v5"
)

func New(wrk worker.Worker) *chi.Mux {
	var (
		r  = chi.NewRouter()
		wr = &workerRoutes{
			wrk: wrk,
		}
	)

	r.Post("/map", wr.handleMap)
	r.Post("/reduce", wr.handleReduce)
	r.Post("/terminate", wr.handleTerminate)

	return r
}

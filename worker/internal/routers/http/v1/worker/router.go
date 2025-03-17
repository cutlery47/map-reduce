package worker

import (
	httpWorker "github.com/cutlery47/map-reduce/worker/internal/domain/http"
	"github.com/go-chi/chi/v5"
)

func New(wrk *httpWorker.Worker, doneChan, recvChan chan<- struct{}, errChan chan<- error) chi.Router {
	var (
		r  = chi.NewRouter()
		wr = &workerRoutes{
			wrk:      wrk,
			doneChan: doneChan,
			errChan:  errChan,
			recvChan: recvChan,
		}
	)

	r.Post("/map", wr.handleMap)
	r.Post("/reduce", wr.handleReduce)

	return r
}

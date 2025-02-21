package worker

import (
	httpService "github.com/cutlery47/map-reduce/worker/internal/service/http"
	"github.com/go-chi/chi/v5"
)

func New(svc *httpService.Service, doneChan, recvChan chan<- struct{}, errChan chan<- error) chi.Router {
	var (
		r  = chi.NewRouter()
		wr = &workerRoutes{
			svc:      svc,
			doneChan: doneChan,
			errChan:  errChan,
			recvChan: recvChan,
		}
	)

	r.Post("/map", wr.handleMap)
	r.Post("/reduce", wr.handleReduce)

	return r
}

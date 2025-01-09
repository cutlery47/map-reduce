package v1

import (
	"net/http"

	"github.com/cutlery47/map-reduce/worker/internal/service"
	"github.com/go-chi/chi/v5"
)

func NewController(r chi.Router, srv service.Service, errChan chan<- error, recvChan chan<- struct{}) {
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("pong"))
	})

	wr := &workerRoutes{
		srv:      srv,
		errChan:  errChan,
		recvChan: recvChan,
	}

	r.Post("/map", wr.handleMap)
	r.Post("/reduce", wr.handleReduce)
}

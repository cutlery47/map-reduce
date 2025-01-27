package controller

import (
	"net/http"

	"github.com/cutlery47/map-reduce/worker/internal/core"
	"github.com/go-chi/chi/v5"
)

func New(r chi.Router, w core.Worker, errChan chan<- error, recvChan, endChan chan<- struct{}) {
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("pong"))
	})

	wr := &workerRoutes{
		w:        w,
		errChan:  errChan,
		endChan:  endChan,
		recvChan: recvChan,
	}

	r.Post("/map", wr.handleMap)
	r.Post("/reduce", wr.handleReduce)
}

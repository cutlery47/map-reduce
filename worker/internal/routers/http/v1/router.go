package handlers

import (
	"net/http"

	httpWorker "github.com/cutlery47/map-reduce/worker/internal/domain/http"
	"github.com/cutlery47/map-reduce/worker/internal/routers/http/v1/worker"
	"github.com/go-chi/chi/v5"
)

func New(
	wrk *httpWorker.Worker,
	doneChan, recvChan chan<- struct{},
	errChan chan<- error,
) *chi.Mux {
	var (
		r = chi.NewRouter()
	)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("pong"))
	})

	r.Mount("/worker", worker.New(wrk, doneChan, recvChan, errChan))

	return r
}

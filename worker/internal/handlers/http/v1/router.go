package handlers

import (
	"net/http"

	"github.com/cutlery47/map-reduce/worker/internal/handlers/http/v1/worker"
	httpService "github.com/cutlery47/map-reduce/worker/internal/service/http"
	"github.com/go-chi/chi/v5"
)

func New(
	svc *httpService.Service,
	doneChan, recvChan chan<- struct{},
	errChan chan<- error,
) chi.Router {
	var (
		r = chi.NewRouter()
	)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("pong"))
	})

	r.Mount("/worker", worker.New(svc, doneChan, recvChan, errChan))

	return r
}

package v1

import (
	"net/http"

	"github.com/cutlery47/map-reduce/master/internal/service"
	"github.com/go-chi/chi/v5"
)

func NewController(r chi.Router, srv service.Service) {
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("pong"))
	})

	rs := &routes{
		srv: srv,
	}

	r.Post("/register", rs.registerWorkers)
}

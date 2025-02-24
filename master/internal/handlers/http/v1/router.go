package v1

import (
	"net/http"

	"github.com/cutlery47/map-reduce/master/internal/handlers/http/v1/register"
	"github.com/cutlery47/map-reduce/master/internal/handlers/http/v1/result"
	"github.com/cutlery47/map-reduce/master/internal/service"
	"github.com/go-chi/chi/v5"
)

func New(srv *service.Master, regChan <-chan bool) *chi.Mux {
	var (
		mux = chi.NewMux()
	)

	mux.Route("/api/v1", func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte("pong"))
		})
		r.Mount("/result", result.New(srv))
		r.Mount("/register", register.New(srv, regChan))

	})

	return mux
}

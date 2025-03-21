package routers

import (
	"net/http"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/master/internal/domain/master"
	"github.com/cutlery47/map-reduce/master/internal/routers/http/v1/register"
	"github.com/cutlery47/map-reduce/master/internal/routers/http/v1/result"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(conf mr.Config, mst *master.Master) *chi.Mux {
	var (
		mux = chi.NewMux()
	)

	mux.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Logger)

		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte("pong"))
		})
		r.Mount("/result", result.New(mst))
		r.Mount("/register", register.New(conf, mst))

	})

	return mux
}

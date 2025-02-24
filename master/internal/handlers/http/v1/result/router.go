package result

import (
	"github.com/cutlery47/map-reduce/master/internal/service"
	"github.com/go-chi/chi/v5"
)

func New(srv *service.Master) *chi.Mux {
	var (
		mux = chi.NewMux()
		rr  = &resultRoutes{
			srv: srv,
		}
	)

	mux.Group(func(r chi.Router) {
		r.Post("/map", rr.mapRes)
		r.Post("/reduce", rr.redRes)
	})

	return mux
}

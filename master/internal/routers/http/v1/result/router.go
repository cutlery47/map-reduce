package result

import (
	"github.com/cutlery47/map-reduce/master/internal/domain/master"
	"github.com/go-chi/chi/v5"
)

func New(mst *master.Master) *chi.Mux {
	var (
		mux = chi.NewMux()
		rr  = &resultRoutes{
			srv: mst,
		}
	)

	mux.Group(func(r chi.Router) {
		r.Post("/map", rr.mapRes)
		r.Post("/reduce", rr.redRes)
	})

	return mux
}

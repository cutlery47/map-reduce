package register

import (
	mr "github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/master/internal/domain/master"
	"github.com/go-chi/chi/v5"
)

func New(conf mr.Config, mst *master.Master) *chi.Mux {
	var (
		mux = chi.NewMux()
		rr  = registerRoutes{
			mst:  mst,
			conf: conf,
		}
	)

	mux.Post("/", rr.register)

	return mux
}

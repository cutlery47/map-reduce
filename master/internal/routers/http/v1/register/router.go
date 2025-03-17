package register

import (
	mr "github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/master/internal/domain"
	"github.com/go-chi/chi/v5"
)

func New(conf mr.Config, mst *domain.Master) *chi.Mux {
	var (
		mux = chi.NewMux()
		rr  = registerRoutes{
			mst:  mst,
			conf: conf,
		}
	)

	mux.Post("/register", rr.register)

	return mux
}

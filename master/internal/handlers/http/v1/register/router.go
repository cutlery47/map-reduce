package register

import (
	"github.com/cutlery47/map-reduce/master/internal/service"
	"github.com/go-chi/chi/v5"
)

func New(srv *service.Master, regChan <-chan bool) *chi.Mux {
	var (
		mux = chi.NewMux()
		rr  = registerRoutes{
			srv:     srv,
			regChan: regChan,
		}
	)

	mux.Get("/", rr.register)

	return mux
}

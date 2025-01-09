package app

import (
	"context"
	"errors"
	"net/http"
	"syscall"

	"github.com/cutlery47/map-reduce/mapreduce"
	v1 "github.com/cutlery47/map-reduce/master/internal/controller/http/v1"
	"github.com/cutlery47/map-reduce/master/internal/service"
	"github.com/cutlery47/map-reduce/master/pkg/httpserver"
	"github.com/go-chi/chi/v5"
)

func Run() error {
	// drop permission limit
	syscall.Umask(0)

	ctx := context.Background()
	conf := mapreduce.KubernetesConfig

	if conf.Mappers != conf.Reducers {
		return errors.New("amount of mappers should be equal to the amount of reducers")
	}

	if conf.Mappers == 0 {
		return errors.New("mappers amount should be > 0")
	}

	if conf.Reducers == 0 {
		return errors.New("reducers amount should be > 0")
	}

	errChan := make(chan error)
	regChan := make(chan bool, conf.Mappers+conf.Reducers)

	srv := service.NewMasterService(conf, http.DefaultClient)

	r := chi.NewRouter()
	v1.NewController(r, srv, regChan)

	go srv.HandleWorkers(errChan, regChan)
	return httpserver.New(r).Run(ctx, errChan)
}

package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"syscall"

	"github.com/cutlery47/map-reduce/mapreduce"
	v1 "github.com/cutlery47/map-reduce/master/internal/controller/http/v1"
	"github.com/cutlery47/map-reduce/master/internal/service"
	"github.com/cutlery47/map-reduce/master/pkg/httpserver"
	"github.com/go-chi/chi/v5"
)

var envLocation = flag.String("env", ".env", "specify env-file name and location")

func Run() error {
	// drop permission limit
	syscall.Umask(0)

	flag.Parse()
	ctx := context.Background()

	conf, err := mapreduce.NewMasterConfig(*envLocation)
	if err != nil {
		return fmt.Errorf("mareduce.NewMasterConfig: %v", err)
	}

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

	srv, err := service.NewMasterService(conf, http.DefaultClient)
	if err != nil {
		return fmt.Errorf("service.NewMasterService: %v", err)
	}

	r := chi.NewRouter()
	v1.NewController(r, srv, regChan)

	go srv.HandleWorkers(errChan, regChan)
	return httpserver.New(r).Run(ctx, errChan)
}

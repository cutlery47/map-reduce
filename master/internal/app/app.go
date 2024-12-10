package app

import (
	"context"
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
	errChan := make(chan error)

	srv := service.NewWorkerService(*mapreduce.DefaultConfig, http.DefaultClient, errChan)

	r := chi.NewRouter()
	v1.NewController(r, srv)

	return httpserver.New(r).Run(ctx, errChan)
}

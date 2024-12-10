package app

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/cutlery47/map-reduce/mapreduce"
	v1 "github.com/cutlery47/map-reduce/worker/internal/controller/http/v1"
	"github.com/cutlery47/map-reduce/worker/internal/service"
	"github.com/cutlery47/map-reduce/worker/pkg/httpserver"
	"github.com/go-chi/chi/v5"
)

func Run() error {
	ctx := context.Background()
	readyChan := make(chan struct{})
	errChan := make(chan error)

	cl := http.DefaultClient
	cl.Timeout = 3 * time.Second

	srv := service.NewWorkerService(cl, *mapreduce.DefaultConfig)

	r := chi.NewRouter()
	v1.NewController(r, srv, errChan)

	// running http server on random available port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("net.Listen: ", err)
	}

	go srv.SendRegister(ctx, listener.Addr().(*net.TCPAddr).Port, readyChan, errChan)
	return httpserver.New(r).Run(ctx, listener, readyChan, errChan)
}

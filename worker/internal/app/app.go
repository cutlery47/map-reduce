package app

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"

	"github.com/cutlery47/map-reduce/mapreduce"
	v1 "github.com/cutlery47/map-reduce/worker/internal/controller/http/v1"
	"github.com/cutlery47/map-reduce/worker/internal/service"
	"github.com/cutlery47/map-reduce/worker/pkg/httpserver"
	"github.com/go-chi/chi/v5"
)

var envLocation = flag.String("env", ".env", "specify env-file name and location")

func Run() error {
	flag.Parse()
	ctx := context.Background()

	// channel for signaling that http server has been set up
	readyChan := make(chan struct{})
	// channel for passing errors
	errChan := make(chan error)
	// channel for ending worker
	endChan := make(chan struct{})
	// channel for signaling that a new job was received
	recvChan := make(chan struct{})

	conf, err := mapreduce.NewConfig(*envLocation)
	if err != nil {
		return fmt.Errorf("error when reading config: %v", err)
	}

	srv := service.NewWorkerService(http.DefaultClient, conf, endChan, recvChan)

	r := chi.NewRouter()
	v1.NewController(r, srv, errChan, recvChan)

	// running http server on random available port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return fmt.Errorf("net.Listen: %v", err)
	}

	go srv.SendRegister(ctx, listener.Addr().(*net.TCPAddr).Port, readyChan, errChan)
	return httpserver.New(r).Run(ctx, listener, readyChan, errChan, endChan)
}

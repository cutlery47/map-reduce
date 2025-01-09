package app

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/cutlery47/map-reduce/mapreduce"
	v1 "github.com/cutlery47/map-reduce/worker/internal/controller/http/v1"
	"github.com/cutlery47/map-reduce/worker/internal/service"
	"github.com/cutlery47/map-reduce/worker/pkg/httpserver"
	"github.com/go-chi/chi/v5"
)

func Run() error {
	ctx := context.Background()

	// channel for signaling that http server has been set up
	readyChan := make(chan struct{})
	// channel for passing errors
	errChan := make(chan error)
	// channel for ending worker
	endChan := make(chan struct{})
	// channel for signaling that a new job was received
	recvChan := make(chan struct{})

	srv := service.NewWorkerService(http.DefaultClient, mapreduce.KubernetesConfig, endChan, recvChan)

	r := chi.NewRouter()
	v1.NewController(r, srv, errChan, recvChan)

	// running http server on random available port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("net.Listen: ", err)
	}

	go srv.SendRegister(ctx, listener.Addr().(*net.TCPAddr).Port, readyChan, errChan)
	return httpserver.New(r).Run(ctx, listener, readyChan, errChan, endChan)
}

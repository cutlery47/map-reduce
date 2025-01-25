package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/worker/internal/core"
	httpworker "github.com/cutlery47/map-reduce/worker/internal/http"
	"github.com/cutlery47/map-reduce/worker/pkg/httpserver"
	"github.com/go-chi/chi/v5"
)

var envLocation = flag.String("env", ".env", "specify env-file name and location")

func Run() error {
	flag.Parse()

	conf, err := mapreduce.NewWorkerConfig(*envLocation)
	if err != nil {
		return fmt.Errorf("error when reading config: %v", err)
	}

	// channel for passing errors
	errChan := make(chan error)

	// worker for handing mapping / reducing
	dw := core.NewDefaultWorker(conf)
	// registrar for announcing worker to the master
	dr := core.NewDefaultRegistrar(http.DefaultClient, errChan, conf.WorkerRegistrarConfig)

	// run preferred worker based on producer type
	switch conf.ProducerType {
	case "HTTP":
		return runHttp(dw, dr, errChan, conf)
	case "QUEUE":
		return runQueue(dw, dr, errChan, conf)
	default:
		return errors.New("undefined producer type")
	}
}

// run http-based worker
func runHttp(w core.Worker, r core.Registrar, errChan chan error, conf mapreduce.WorkerConfig) error {
	ctx := context.Background()

	// channel for signaling that http server has been set up
	readyChan := make(chan struct{})
	// channel for signaling that a new job was received
	recvChan := make(chan struct{})
	// channel for signaling that worker has ended execution
	endChan := make(chan struct{})

	// running http server on random available port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return fmt.Errorf("net.Listen: %v", err)
	}

	// creating http-controller for receiving tasks from master
	rt := chi.NewRouter()
	httpworker.NewController(rt, w, errChan, recvChan, endChan)

	// announcing worker to master
	body := mapreduce.WorkerRegisterRequest{Port: strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)}
	send := httpworker.UsingHttpWorker(ctx, body, r.SendRegister)
	go send(conf.SetupDuration, conf.WorkerTimeout, errChan, readyChan, recvChan)

	// running http server
	return httpserver.New(rt, listener, readyChan, errChan, endChan).Run(ctx)
}

// run queue-based worker
func runQueue(w core.Worker, r core.Registrar, errChan chan error, conf mapreduce.WorkerConfig) error {
	ctx := context.Background()

	go r.SendRegister(ctx, mapreduce.WorkerRegisterRequest{Port: "1337"})

	return <-errChan
}

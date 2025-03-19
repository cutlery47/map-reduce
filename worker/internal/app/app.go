package app

import (
	"errors"
	"fmt"
	"net"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/worker/internal/domain/worker"
	httpWorker "github.com/cutlery47/map-reduce/worker/internal/domain/worker/http"
	rabbitWorker "github.com/cutlery47/map-reduce/worker/internal/domain/worker/rabbit"
	routers "github.com/cutlery47/map-reduce/worker/internal/routers/http/v1"
	hserv "github.com/cutlery47/map-reduce/worker/pkg/httpserver"
	log "github.com/sirupsen/logrus"
)

func Run(conf mr.Config) error {
	log.Infoln("[SETUP] setting up worker...")

	// running socket listener on random available port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return err
	}

	var (
		wrk    worker.Worker                         // worker instance
		doneCh = make(chan hserv.AppSignal)          // channel for passing possible errors down to httpserver
		port   = listener.Addr().(*net.TCPAddr).Port // port on which worker is running
	)

	// creating worker based on transport
	switch conf.Transport {
	case "HTTP":
		httpWrk, err := httpWorker.New(conf, port)
		if err != nil {
			return fmt.Errorf("error setting up http worker: %v", err)
		}
		wrk = httpWrk
	case "QUEUE":
		rabbitWrk, err := rabbitWorker.New(conf)
		if err != nil {
			return fmt.Errorf("error setting up rabbitmq worker: %v", err)
		}
		wrk = rabbitWrk
	default:
		return errors.New("undefined transport")
	}

	go func() {
		var sig hserv.AppSignal

		err := wrk.Run()
		if err != nil {
			sig.Error = err
		} else {
			sig.Message = "Success"
		}

		doneCh <- sig
	}()

	// creating http-controller for receiving tasks from master
	rt := routers.New(wrk)

	// running http server
	return hserv.New(rt, listener,
		hserv.WithReadTimeout(conf.WorkerReadTimeout),
		hserv.WithWriteTimeout(conf.WorkerWriteTimeout),
		hserv.WithShutdownTimeout(conf.WorkerShutdownTimeout),
	).Run(doneCh)
}

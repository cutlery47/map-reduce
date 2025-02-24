package app

import (
	"errors"
	"flag"
	"fmt"
	"log"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/worker/internal/core"
	v1 "github.com/cutlery47/map-reduce/worker/internal/handlers/http/v1"
	httpService "github.com/cutlery47/map-reduce/worker/internal/service/http"
	queueService "github.com/cutlery47/map-reduce/worker/internal/service/queue"
	"github.com/cutlery47/map-reduce/worker/pkg/httpserver"
)

var envLocation = flag.String("env", ".env", "specify env-file name and location")

func Run() error {
	flag.Parse()

	conf, err := mr.NewWorkerConfig(*envLocation)
	if err != nil {
		return fmt.Errorf("error when reading config: %v", err)
	}

	var (
		// creating core worker for handing mapping / reducing
		w = core.NewWorker(conf.WrkCoreConf)
		// creaing registrar for announcing worker to the master
		r = core.NewRegistrar(conf.WrkRegConf)
	)

	// creating service based on transport
	switch conf.Transport {
	case "HTTP":
		return runHttp(conf.WrkSvcConf, w, r)
	case "QUEUE":
		return runQueue(conf.WrkSvcConf, w, r)
	default:
		return errors.New("undefined transport")
	}
}

// run http-based worker
func runHttp(conf mr.WrkSvcConf, w *core.Worker, r *core.Registrar) error {
	svc, err := httpService.New(conf, w, r)
	if err != nil {
		return fmt.Errorf("error when creating http service: %v", err)
	}

	// running http server on random available port
	var (
		// channel for signaling that worker has finished execution
		doneChan = make(chan struct{})
		// channel for passing errors down to httpserver
		errChan = make(chan error)
		// channel for signaling that a new job was received
		recvChan = make(chan struct{})
		// channel for signaling that http server has been set up
		readyChan = make(chan struct{})
	)

	// creating http-controller for receiving tasks from master
	h := v1.New(svc, doneChan, recvChan, errChan)

	go func() {
		err := svc.Run(readyChan, recvChan)
		if err != nil {
			// sending all runtime errors to httpserver
			errChan <- err
		}
	}()

	// running http server
	return httpserver.New(conf.WrkHttpConf, h).Run(doneChan, readyChan, errChan)
}

// run queue-based worker
func runQueue(conf mr.WrkSvcConf, w *core.Worker, r *core.Registrar) error {
	// br, err := queue.NewBrocker(conf.RabbitConf)
	// if err != nil {
	// 	return fmt.Errorf("error when creating new brocker: %v", err)
	// }

	svc, err := queueService.New(conf, w, r)
	if err != nil {
		return fmt.Errorf("error when creating queue service: %v", err)
	}

	var (
		// channel for signaling that worker has finished execution
		doneChan chan struct{}
		// channel for passing errors down to httpserver
		errChan chan error
	)

	go func() {
		err := svc.Run(doneChan)
		if err != nil {
			// sending all runtime errors to httpserver
			errChan <- err
		}
	}()

	select {
	case <-doneChan:
		log.Println("service is done")
		return nil
	case err := <-errChan:
		return err
	}
}

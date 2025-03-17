package app

import (
	"errors"
	"net"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	httpWorker "github.com/cutlery47/map-reduce/worker/internal/domain/http"
	rabbitWorker "github.com/cutlery47/map-reduce/worker/internal/domain/rabbit"
	v1 "github.com/cutlery47/map-reduce/worker/internal/routers/http/v1"
	"github.com/cutlery47/map-reduce/worker/pkg/httpserver"
	log "github.com/sirupsen/logrus"
)

func Run(conf mr.Config) error {
	log.Infoln("[SETUP] setting up worker...")

	// creating service based on transport
	switch conf.Transport {
	case "HTTP":
		return runHttp(conf)
	case "QUEUE":
		return runQueue(conf)
	default:
		return errors.New("undefined transport")
	}
}

// run http-based worker
func runHttp(conf mr.Config) error {
	wrk, err := httpWorker.New(conf)
	if err != nil {
		return err
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
	rt := v1.New(wrk, doneChan, recvChan, errChan)

	go func() {
		err := wrk.Run(readyChan, recvChan)
		if err != nil {
			// sending all runtime errors to httpserver
			errChan <- err
		}
	}()

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return err
	}

	// running http server
	return httpserver.New(conf, listener, rt).Run(doneChan, readyChan, errChan)
}

// run queue-based worker
func runQueue(conf mr.Config) error {
	wrk, err := rabbitWorker.New(conf)
	if err != nil {
		return err
	}

	var (
		// channel for signaling that worker has finished execution
		doneChan chan struct{}
		// channel for passing errors down to httpserver
		errChan chan error
	)

	go func() {
		err := wrk.Run(doneChan)
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

package app

import (
	"fmt"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/master/internal/domain/master"
	prod "github.com/cutlery47/map-reduce/master/internal/domain/producer"
	httpProducer "github.com/cutlery47/map-reduce/master/internal/domain/producer/http"
	rabbitProducer "github.com/cutlery47/map-reduce/master/internal/domain/producer/rabbit"
	routers "github.com/cutlery47/map-reduce/master/internal/routers/http/v1"
	"github.com/cutlery47/map-reduce/master/pkg/httpserver"
)

func Run(conf mr.Config) error {
	var (
		tp prod.TaskProducer
	)

	// creating task producer based on config
	switch conf.Transport {
	case "HTTP":
		tp = httpProducer.New(conf)
	case "QUEUE":
		rtp, err := rabbitProducer.New(conf)
		if err != nil {
			return fmt.Errorf("[SETUP] error when setting up rabbitmq task producer: %v", err)
		}
		tp = rtp
	}

	mst, err := master.New(conf, tp)
	if err != nil {
		return fmt.Errorf("[SETUP] error when setting up master: %v", err)
	}

	rt := routers.New(conf, mst)

	var (
		doneChan = make(chan struct{})
		errChan  = make(chan error)
	)

	// entrypoint
	go func() {
		err := mst.Run()
		if err != nil {
			// error occured => pass to httpserver
			errChan <- err
		} else {
			// all fine
			doneChan <- struct{}{}
		}
	}()

	return httpserver.New(conf, rt).Run(doneChan, errChan)
}

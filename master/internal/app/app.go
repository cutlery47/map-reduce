package app

import (
	"fmt"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/master/internal/domain"
	prod "github.com/cutlery47/map-reduce/master/internal/domain/producer"
	httpProducer "github.com/cutlery47/map-reduce/master/internal/domain/producer/http"
	rabbitProducer "github.com/cutlery47/map-reduce/master/internal/domain/producer/rabbit"
	routers "github.com/cutlery47/map-reduce/master/internal/routers/http/v1"
	"github.com/cutlery47/map-reduce/master/pkg/httpserver"
)

func Run(conf mr.Config) error {

	var (
		doneChan = make(chan struct{})
		errChan  = make(chan error)
	)

	var (
		tp prod.TaskProducer
	)

	switch conf.Transport {
	case "http":
		tp = httpProducer.New(conf)
	case "queue":
		rtp, err := rabbitProducer.New(conf)
		if err != nil {
			return fmt.Errorf("[SETUP] error when setting up rabbitmq task producer: %v", err)
		}

		tp = rtp
	}

	var (
		mst = domain.NewMaster(conf, tp)
		rt  = routers.New(conf, mst)
	)

	// entrypoint
	go func() {
		err := mst.Work()
		if err != nil {
			errChan <- err
		} else {
			doneChan <- struct{}{}
		}
	}()

	return httpserver.New(conf, rt).Run(doneChan, errChan)
}

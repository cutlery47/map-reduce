package app

import (
	"errors"
	"fmt"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/master/internal/domain/master"
	prod "github.com/cutlery47/map-reduce/master/internal/domain/producer"
	httpProducer "github.com/cutlery47/map-reduce/master/internal/domain/producer/http"
	rabbitProducer "github.com/cutlery47/map-reduce/master/internal/domain/producer/rabbit"
	routers "github.com/cutlery47/map-reduce/master/internal/routers/http/v1"
	hserv "github.com/cutlery47/map-reduce/master/pkg/httpserver"
	log "github.com/sirupsen/logrus"
)

func Run(conf mr.Config) error {
	log.Infoln("[SETUP] setting up master...")

	if conf.Mappers == 0 || conf.Reducers == 0 {
		return errors.New("worker number must be positive")
	}

	var (
		tp     prod.TaskProducer            // task producer instance
		doneCh = make(chan hserv.AppSignal) // channel for passing signals down to httpserver
	)

	// creating task producer based on transport
	switch conf.Transport {
	case "HTTP":
		tp = httpProducer.New(conf)
	case "QUEUE":
		rtp, err := rabbitProducer.New(conf)
		if err != nil {
			return fmt.Errorf("error etting up rabbitmq task producer: %v", err)
		}
		tp = rtp
	default:
		return errors.New("undefined transport")
	}

	// creating master instance
	mst, err := master.New(conf, tp)
	if err != nil {
		return fmt.Errorf("error setting up master: %v", err)
	}

	// running master
	go func() {
		var sig hserv.AppSignal

		err := mst.Run()
		if err != nil {
			sig.Error = err
		} else {
			sig.Message = "Success"
		}

		doneCh <- sig
	}()

	var (
		rt   = routers.New(conf, mst)                                 // router for receiving tasks from master
		addr = fmt.Sprintf("%v:%v", conf.MasterHost, conf.MasterPort) // bind addr
	)

	// running http server
	return hserv.New(rt,
		hserv.WithAddr(addr),
		hserv.WithReadTimeout(conf.MasterReadTimeout),
		hserv.WithWriteTimeout(conf.MasterWriteTimeout),
		hserv.WithShutdownTimeout(conf.MasterShutdownTimeout),
	).Run(doneCh)
}

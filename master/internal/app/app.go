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
	"github.com/cutlery47/map-reduce/master/pkg/httpserver"
	log "github.com/sirupsen/logrus"
)

func Run(conf mr.Config) error {
	log.Infoln("[SETUP] setting up master...")

	var (
		tp     prod.TaskProducer  // task producer instance
		doneCh = make(chan error) // channel for passing possible errors down to httpserver
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

	go func() {
		doneCh <- mst.Run()
	}()

	var (
		rt   = routers.New(conf, mst)                                 // router for receiving tasks from master
		addr = fmt.Sprintf("%v:%v", conf.MasterHost, conf.MasterPort) // bind addr
	)

	// running http server
	return httpserver.New(rt,
		httpserver.WithAddr(addr),
		httpserver.WithReadTimeout(conf.MasterReadTimeout),
		httpserver.WithWriteTimeout(conf.MasterWriteTimeout),
		httpserver.WithShutdownTimeout(conf.MasterShutdownTimeout),
	).Run(doneCh)
}

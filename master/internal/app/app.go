package app

import (
	"flag"
	"fmt"
	"syscall"

	"github.com/cutlery47/map-reduce/mapreduce"
	v1 "github.com/cutlery47/map-reduce/master/internal/handlers/http/v1"
	"github.com/cutlery47/map-reduce/master/pkg/httpserver"
)

var envLocation = flag.String("env", ".env", "specify env-file name and location")

func Run() error {
	// drop permission limit
	syscall.Umask(0)

	flag.Parse()

	conf, err := mapreduce.NewMasterConfig(*envLocation)
	if err != nil {
		return fmt.Errorf("mareduce.NewMasterConfig: %v", err)
	}

	var (
		doneChan = make(chan struct{})
		errChan  = make(chan error)
		regChan  = make(chan bool, conf.Mappers+conf.Reducers)
	)

	h := v1.New(srv, regChan)

	go srv.HandleWorkers(errChan, regChan)
	return httpserver.New(conf.MstHttpConf, h).Run(errChan)
}

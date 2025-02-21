package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/worker/internal/core"
)

type Service struct {
	w *core.Worker
	r *core.Registrar

	cl *http.Client

	conf mr.WrkSvcConf
}

func New(conf mr.WrkSvcConf, w *core.Worker, r *core.Registrar) (*Service, error) {
	return &Service{
		w:  w,
		r:  r,
		cl: http.DefaultClient,
	}, nil
}

func (s *Service) Run(readyChan, recvChan <-chan struct{}) error {
	// waiting for http server to set up
	readyTimer := time.NewTimer(s.conf.SetupDuration)
	select {
	case <-readyTimer.C:
		return fmt.Errorf("server setup timed out")
	case <-readyChan:
	}

	ctx, cancel := context.WithTimeout(context.TODO(), s.conf.RegisterTimeout)
	defer cancel()

	body := mr.WrkRegReq{
		Host: s.conf.WrkHttpConf.Host,
		Port: s.conf.WrkHttpConf.Port,
	}

	// registering on master
	s.r.Send(ctx, body)

	// waiting for new jobs
	// terminate worker if no jobs were received
	timer := time.NewTimer(s.conf.WorkerTimeout)
	log.Println("waiting for tasks for:", s.conf.WorkerTimeout)

	select {
	case <-recvChan:
		return nil
	case <-timer.C:
		return fmt.Errorf("didn't receive any tasks")
	}
}

func (s *Service) Map(r io.Reader) (any, error) {
	return s.w.Map(r)
}

func (s *Service) Reduce(mapped any) (any, error) {
	return s.w.Reduce(mapped)
}

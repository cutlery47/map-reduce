package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/worker/internal/domain/core"
	log "github.com/sirupsen/logrus"
)

type Worker struct {
	mrh  *core.MapReduceHandler
	reg  *core.RegisterHandler
	cl   *http.Client
	conf mr.Config
}

func New(conf mr.Config) (*Worker, error) {
	mrh, err := core.NewMapReduceHandler(conf)
	if err != nil {
		return nil, err
	}

	reg, err := core.NewRegisterHandler(conf)
	if err != nil {
		return nil, err
	}

	return &Worker{
		mrh:  mrh,
		reg:  reg,
		cl:   http.DefaultClient,
		conf: conf,
	}, nil
}

func (w *Worker) Run(port string, readyChan, recvChan <-chan struct{}) error {
	// waiting for http server to set up
	readyTimer := time.NewTimer(w.conf.WorkerSetupDur)
	select {
	case <-readyTimer.C:
		return fmt.Errorf("server setup timed out")
	case <-readyChan:
	}

	ctx, cancel := context.WithTimeout(context.Background(), w.conf.WorkerAwaitDur)
	defer cancel()

	body := mr.RegisterRequest{
		Addr: mr.Addr{
			Host: w.conf.WorkerHost,
			Port: port,
		},
	}

	// registering on master
	w.reg.Send(ctx, body)

	// waiting for new jobs
	// terminate worker if no jobs were received
	timer := time.NewTimer(w.conf.WorkerAwaitDur)
	log.Infoln("waiting for tasks for:", w.conf.WorkerAwaitDur)

	select {
	case <-recvChan:
		return nil
	case <-timer.C:
		return fmt.Errorf("didn't receive any tasks")
	}
}

func (w *Worker) Map(r io.Reader) (any, error) {
	return w.mrh.Map(r)
}

func (w *Worker) Reduce(mapped any) (any, error) {
	return w.mrh.Reduce(mapped)
}

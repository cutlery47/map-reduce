package httpworker

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/worker/internal/domain/core"
	log "github.com/sirupsen/logrus"
)

type HTTPWorker struct {
	reg *core.RegisterHandler
	cl  *http.Client

	recvCh chan struct{} // channel for signaling that a task was received
	port   int
	conf   mr.Config
}

func New(conf mr.Config, port int) (*HTTPWorker, error) {
	reg, err := core.NewRegisterHandler(conf)
	if err != nil {
		return nil, err
	}

	return &HTTPWorker{
		reg:    reg,
		cl:     http.DefaultClient,
		recvCh: make(chan struct{}),
		port:   port,
		conf:   conf,
	}, nil
}

func (w *HTTPWorker) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), w.conf.WorkerAwaitDur)
	defer cancel()

	body := mr.RegisterRequest{
		Addr: mr.Addr{
			Host: w.conf.WorkerHost,
			Port: strconv.Itoa(w.port),
		},
	}

	// registering on master
	_, err := w.reg.Send(ctx, body)
	if err != nil {
		return err
	}

	// waiting for new jobs
	// terminate worker if no jobs were received
	timer := time.NewTimer(w.conf.WorkerAwaitDur)
	log.Infoln("waiting for tasks for:", w.conf.WorkerAwaitDur)

	select {
	case <-w.recvCh:
		return nil
	case <-timer.C:
		return fmt.Errorf("didn't receive any tasks")
	}
}

func (w *HTTPWorker) Map(input io.Reader) (io.Reader, error) {
	return mr.MapperFunc(input)
}

func (w *HTTPWorker) Reduce(input io.Reader) (io.Reader, error) {
	return mr.ReducerFunc(input)
}

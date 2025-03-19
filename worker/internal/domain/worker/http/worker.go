package httpworker

import (
	"errors"
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

	role *mr.Role
	port int
	conf mr.Config

	recvCh chan struct{}            // channel for signaling that a job has been received
	doneCh chan error               // channel for signaling that a job has been handled
	termCh chan mr.TerminateMessage // channel for signaling that worker must terminate
}

func New(conf mr.Config, port int) (*HTTPWorker, error) {
	reg, err := core.NewRegisterHandler(conf)
	if err != nil {
		return nil, err
	}

	return &HTTPWorker{
		reg:    reg,
		cl:     http.DefaultClient,
		port:   port,
		conf:   conf,
		recvCh: make(chan struct{}),
		doneCh: make(chan error),
		termCh: make(chan mr.TerminateMessage),
	}, nil
}

func (w *HTTPWorker) Run() error {
	// registering on master
	res, err := w.reg.Register(mr.Addr{
		Host: w.conf.WorkerHost,
		Port: strconv.Itoa(w.port),
	})
	if err != nil {
		return err
	}

	if res == nil {
		// this shouldn't happen, however just in case
		return errors.New("received nil registration response")
	}

	// assigning role from response
	w.role = &res.Role

	// receive jobs from master until theres no more
	//
	// on each received job the job timer is reset
	// if no jobs were sent during the job timer duration -> terminate
	// if terminate request was sent by master -> terminate

	var (
		timer = time.NewTimer(w.conf.RegisterDur + w.conf.WorkerAwaitDur)
	)

	for {
		log.Infoln("[WORKER] waiting for jobs...")

		select {
		case <-w.recvCh:
			log.Infoln("[WORKER] received a job")
			// proceed to job handling
		case msg := <-w.termCh:
			return fmt.Errorf("received terminate message: %v, with code: %v", msg.Description, msg.Code)
		case <-timer.C:
			return errors.New("no jobs have been received")
		}

		err := <-w.doneCh
		if err != nil {
			return fmt.Errorf("error when handling a job: %v", err)
		}
		log.Infoln("[WORKER] job handled successfully")

		timer.Reset(w.conf.WorkerAwaitDur)
	}

}

// passes job to map-function handler
func (w *HTTPWorker) Map(input io.Reader) (io.Reader, error) {
	w.recvCh <- struct{}{} // job has been received
	var (
		res, err = mr.MapperFunc(input)
	)
	w.doneCh <- err // job has been handled

	return res, err
}

// passes job to reduce-function handler
func (w *HTTPWorker) Reduce(input io.Reader) (io.Reader, error) {
	w.recvCh <- struct{}{} // job has been received
	var (
		res, err = mr.ReducerFunc(input)
	)
	w.doneCh <- err // job has been handled

	return res, err
}

// tells worker to stop execution
func (w *HTTPWorker) Terminate(msg mr.TerminateMessage) {
	w.termCh <- msg
}

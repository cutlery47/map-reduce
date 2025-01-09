package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/cutlery47/map-reduce/mapreduce"
)

var mapFunc, reduceFunc = mapreduce.MapReduce()

type Service interface {
	Map(reader io.Reader) (interface{}, error)
	Reduce(mapped interface{}) (interface{}, error)
}

type WorkerService struct {
	// client for contacting master node
	cl *http.Client

	// channel for signaling that the worker has finished execution
	endChan chan<- struct{}
	// channel for checking that worker has received a job from master
	recvChan <-chan struct{}

	// config
	conf mapreduce.Config
}

func NewWorkerService(cl *http.Client, conf mapreduce.Config, endChan chan<- struct{}, recvChan <-chan struct{}) *WorkerService {
	return &WorkerService{
		cl:       cl,
		conf:     conf,
		endChan:  endChan,
		recvChan: recvChan,
	}
}

func (ws *WorkerService) Map(reader io.Reader) (interface{}, error) {
	log.Println("in map")

	mapResult, err := mapFunc(reader)
	if err != nil {
		return nil, err
	}

	log.Println("finished map")
	ws.endChan <- struct{}{}

	return mapResult, nil
}

func (ws *WorkerService) Reduce(result interface{}) (interface{}, error) {
	log.Println("in reduce")

	mapResult, ok := result.(mapreduce.MyMapResult)
	if !ok {
		return nil, fmt.Errorf("map result incompatible with reduce func")
	}

	redResult, err := reduceFunc(mapResult)
	if err != nil {
		return nil, err
	}

	log.Println("finished reduce")
	ws.endChan <- struct{}{}

	return redResult, nil
}

// Registering on master node
func (ws *WorkerService) SendRegister(ctx context.Context, port int, readyChan <-chan struct{}, errChan chan<- error) {
	// waiting for server setup
	readyTimer := time.NewTimer(ws.conf.ReadyTimeout)
	select {
	case <-readyTimer.C:
		errChan <- fmt.Errorf("server setup timed out")
		return
	case <-readyChan:
	}

	body := mapreduce.WorkerRegisterRequest{
		Port: strconv.Itoa(port),
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		errChan <- fmt.Errorf("json.Marshal: %v", err)
		return
	}
	jsonReader := bytes.NewReader(jsonBody)

	// creating registration request
	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:8080/register", jsonReader)
	if err != nil {
		errChan <- fmt.Errorf("http.NewRequestWithContext: %v", err)
		return
	}

	// sending req
	res, err := ws.cl.Do(req)
	if err != nil {
		errChan <- fmt.Errorf("ws.cl.Do: %v", err)
		return
	}

	if res.StatusCode == 500 {
		errChan <- fmt.Errorf("couldn't register on master node")
		return
	}

	// waiting for new jobs
	// terminate worker if no jobs were received
	go func() {
		timer := time.NewTimer(ws.conf.WorkerTimeout)

		select {
		case <-ws.recvChan:
			return
		case <-timer.C:
			errChan <- fmt.Errorf("didn't receive any jobs")
			return
		}
	}()
}

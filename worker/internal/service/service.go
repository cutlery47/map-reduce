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

	conf mapreduce.Config
}

func NewWorkerService(cl *http.Client, conf mapreduce.Config) *WorkerService {
	return &WorkerService{
		cl:   cl,
		conf: conf,
	}
}

func (ws *WorkerService) Map(reader io.Reader) (interface{}, error) {
	mapResult, err := mapFunc(reader)
	if err != nil {
		return nil, err
	}

	return mapResult, nil
}

func (ws *WorkerService) Reduce(result interface{}) (interface{}, error) {
	mapResult, ok := result.(mapreduce.MyMapResult)
	if !ok {
		return nil, fmt.Errorf("map result incompatible with reduce func")
	}

	redResult, err := reduceFunc(mapResult)
	if err != nil {
		return nil, err
	}

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

	reqCtx, cancel := context.WithTimeout(ctx, ws.conf.RequestTimeout)
	defer cancel()

	// sending registration request
	req, err := http.NewRequestWithContext(reqCtx, "POST", "http://localhost:8080/register", jsonReader)
	if err != nil {
		errChan <- fmt.Errorf("http.NewRequestWithContext: %v", err)
		return
	}

	res, err := ws.cl.Do(req)
	if err != nil {
		errChan <- fmt.Errorf("http.Get: %v", err)
		return
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		errChan <- fmt.Errorf("io.ReadAll: %v", err)
		return
	}

	log.Println("master resonse:", string(resBody))
}

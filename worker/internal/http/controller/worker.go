package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/worker/internal/core"
)

type workerRoutes struct {
	w core.Worker

	// channel for signaling that the worker has finished execution
	errChan chan<- error
	// channel for signaling that worker has received a job from master
	recvChan chan<- struct{}
	// channel for signaling that worker has finished handing its tasks
	endChan chan<- struct{}
}

func (wr *workerRoutes) handleMap(w http.ResponseWriter, r *http.Request) {
	// signaling that a new job has been received
	wr.recvChan <- struct{}{}

	res, err := wr.w.Map(r.Body)
	if err != nil {
		handleErr(err, w)
		wr.errChan <- err
		return
	}

	json, err := json.Marshal(res)
	if err != nil {
		handleErr(err, w)
		wr.errChan <- err
		return
	}

	w.WriteHeader(200)
	w.Write(json)

	wr.endChan <- struct{}{}
}

func (wr *workerRoutes) handleReduce(w http.ResponseWriter, r *http.Request) {
	// signaling that a new job has been received
	wr.recvChan <- struct{}{}

	bytesBody, err := io.ReadAll(r.Body)
	if err != nil {
		handleErr(err, w)
		wr.errChan <- fmt.Errorf("io.ReadAll: %v", err)
		return
	}

	var mapResult mapreduce.MyMapResult

	if err := json.Unmarshal(bytesBody, &mapResult); err != nil {
		handleErr(err, w)
		wr.errChan <- fmt.Errorf("json.Unmarshall: %v", err)
		return
	}

	res, err := wr.w.Reduce(mapResult)
	if err != nil {
		handleErr(err, w)
		wr.errChan <- fmt.Errorf("reduceFunc: %v", err)
		return
	}

	jsonRes, err := json.Marshal(res)
	if err != nil {
		handleErr(err, w)
		wr.errChan <- fmt.Errorf("json.Marshall: %v", err)
		return
	}

	w.WriteHeader(200)
	w.Write(jsonRes)

	wr.endChan <- struct{}{}
}

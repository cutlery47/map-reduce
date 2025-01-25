package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cutlery47/map-reduce/mapreduce"
)

type SenderFunction func(ctx context.Context, body mapreduce.WorkerRegisterRequest)

type Registrar interface {
	// sender function
	SendRegister(ctx context.Context, body mapreduce.WorkerRegisterRequest)
}

type DefaultRegistrar struct {
	cl *http.Client

	errChan chan<- error

	conf mapreduce.WorkerRegistrarConfig
}

func NewDefaultRegistrar(cl *http.Client, errChan chan<- error, conf mapreduce.WorkerRegistrarConfig) *DefaultRegistrar {
	return &DefaultRegistrar{
		cl:      cl,
		errChan: errChan,
		conf:    conf,
	}
}

func (dr *DefaultRegistrar) SendRegister(ctx context.Context, body mapreduce.WorkerRegisterRequest) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		dr.errChan <- fmt.Errorf("json.Marshal: %v", err)
		return
	}
	jsonReader := bytes.NewReader(jsonBody)

	addr := fmt.Sprintf("http://%v:%v/register", dr.conf.MasterHost, dr.conf.MasterPort)

	// creating registration request
	req, err := http.NewRequestWithContext(ctx, "POST", addr, jsonReader)
	if err != nil {
		dr.errChan <- fmt.Errorf("http.NewRequestWithContext: %v", err)
		return
	}

	// sending req
	res, err := dr.cl.Do(req)
	if err != nil {
		dr.errChan <- fmt.Errorf("ws.cl.Do: %v", err)
		return
	}

	if res.StatusCode != 200 {
		dr.errChan <- fmt.Errorf("couldn't register on master node")
		return
	}
}

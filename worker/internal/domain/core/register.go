package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	mr "github.com/cutlery47/map-reduce/mapreduce"
)

type RegisterHandler struct {
	cl   *http.Client // client for communicating with master
	conf mr.Config
}

func NewRegisterHandler(conf mr.Config) (*RegisterHandler, error) {
	return &RegisterHandler{
		cl:   http.DefaultClient,
		conf: conf,
	}, nil
}

func (rh *RegisterHandler) Send(ctx context.Context, body mr.RegisterRequest) (*mr.RegisterResponse, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %v", err)
	}

	// master addr
	addr := fmt.Sprintf("http://%v:%v/register", rh.conf.MasterHost, rh.conf.MasterPort)

	req, err := http.NewRequestWithContext(ctx, "POST", addr, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %v", err)
	}

	// announcing to master
	res, err := rh.cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ws.cl.Do: %v", err)
	}

	if res.StatusCode != 200 {
		msg, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("couldn't register on master node: %v", string(msg))
	}

	var (
		resJson mr.RegisterResponse
	)

	err = json.NewDecoder(res.Body).Decode(&resJson)
	if err != nil {
		return nil, fmt.Errorf("json.Decode: %v", err)
	}

	return &resJson, nil
}

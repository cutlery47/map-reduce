package core

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
		return nil, err
	}

	// master addr
	addr := fmt.Sprintf("http://%v:%v/register", rh.conf.MasterHost, rh.conf.MasterPort)

	req, err := http.NewRequestWithContext(ctx, "POST", addr, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	// announcing to master
	res, err := rh.cl.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("bad status")
	}

	var (
		resJson mr.RegisterResponse
	)

	err = json.NewDecoder(res.Body).Decode(&resJson)
	if err != nil {
		return nil, err
	}

	return &resJson, nil
}

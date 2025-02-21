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

// default registrar impl
type Registrar struct {
	// for communicating with master
	cl *http.Client

	conf mr.WrkRegConf
}

func NewRegistrar(conf mr.WrkRegConf) *Registrar {
	return &Registrar{
		cl:   http.DefaultClient,
		conf: conf,
	}
}

func (r *Registrar) Send(ctx context.Context, body mr.WrkRegReq) (*mr.WrkRegRes, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %v", err)
	}

	// master addr
	addr := fmt.Sprintf("http://%v:%v/register", r.conf.MasterHost, r.conf.MasterPort)

	req, err := http.NewRequestWithContext(ctx, "POST", addr, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %v", err)
	}

	// announcing to master
	res, err := r.cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ws.cl.Do: %v", err)
	}

	if res.StatusCode != 200 {
		msg, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("couldn't register on master node: %v", string(msg))
	}

	var (
		resJson mr.WrkRegRes
	)

	err = json.NewDecoder(res.Body).Decode(&resJson)
	if err != nil {
		return nil, fmt.Errorf("json.Decode: %v", err)
	}

	return &resJson, nil
}

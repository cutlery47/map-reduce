package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/mapreduce/models"
	"github.com/cutlery47/map-reduce/mapreduce/requests"
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

func (rh *RegisterHandler) Register(addr models.Addr) (*requests.RegisterResponse, error) {
	var (
		masterAddr = fmt.Sprintf("http://%v:%v/api/v1/register/", rh.conf.MasterHost, rh.conf.MasterPort)
		body       = requests.RegisterRequest{
			Addr: addr,
		}
	)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	res, err := rh.cl.Post(masterAddr, "application/json", bytes.NewReader(jsonBody))
	if res.StatusCode != http.StatusOK {
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("registration failed: %v", string(resBody))
	}

	var (
		resJson requests.RegisterResponse
	)

	err = json.NewDecoder(res.Body).Decode(&resJson)
	if err != nil {
		return nil, err
	}

	return &resJson, nil
}

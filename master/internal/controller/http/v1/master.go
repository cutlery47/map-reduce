package v1

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/master/internal/service"
)

type masterRoutes struct {
	// service
	srv service.Service

	// channel for receiving registration status
	regChan <-chan bool
}

func (mr *masterRoutes) registerWorkers(w http.ResponseWriter, r *http.Request) {
	req := mapreduce.WorkerRegisterRequest{
		Host: strings.Split(r.RemoteAddr, ":")[0],
	}

	jsonBody, err := io.ReadAll(r.Body)
	if err != nil {
		handleErr(err, w)
		return
	}

	if err = json.Unmarshal(jsonBody, &req); err != nil {
		handleErr(err, w)
		return
	}

	if err := mr.srv.Register(req); err != nil {
		handleErr(err, w)
		return
	}

	// waiting for all workers to announce themselves
	// if the amount of announced workers is not sufficient - the request is rejected
	// else - accepted
	registered := <-mr.regChan
	if registered {
		w.WriteHeader(200)
		w.Write([]byte("registered"))
	} else {
		w.WriteHeader(500)
		w.Write([]byte("couldn't collect enough workers in time"))
	}
}

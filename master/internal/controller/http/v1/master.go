package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/master/internal/service"
)

type masterRoutes struct {
	srv service.Service

	errChan chan<- error
}

func (mr *masterRoutes) registerWorkers(w http.ResponseWriter, r *http.Request) {
	req := mapreduce.WorkerRegisterRequest{
		Host: strings.Split(r.RemoteAddr, ":")[0],
	}

	jsonBody, err := io.ReadAll(r.Body)
	if err != nil {
		mr.errChan <- fmt.Errorf("io.ReadAll: %v", err)
		handleErr(err, w)
		return
	}

	if err = json.Unmarshal(jsonBody, &req); err != nil {
		mr.errChan <- fmt.Errorf("json.Unmarshall: %v", err)
		handleErr(err, w)
		return
	}

	if err := mr.srv.Register(req); err != nil {
		mr.errChan <- fmt.Errorf("mr.srv.Register: %v", err)
		handleErr(err, w)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("registered"))
}

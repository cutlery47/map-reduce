package v1

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/master/internal/service"
)

type routes struct {
	srv service.Service
}

func (rs *routes) registerWorkers(w http.ResponseWriter, r *http.Request) {
	log.Println("here")

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

	if err := rs.srv.Register(req); err != nil {
		handleErr(err, w)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("registered"))
}

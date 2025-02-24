package register

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/master/internal/service"
)

type registerRoutes struct {
	srv *service.Master
	// channel for receiving registration status
	regChan <-chan bool
}

func (rr *registerRoutes) register(w http.ResponseWriter, r *http.Request) {
	var req mr.WrkRegReq

	jsonBody, err := io.ReadAll(r.Body)
	if err != nil {
		handleErr(err, w)
		return
	}

	if err = json.Unmarshal(jsonBody, &req); err != nil {
		handleErr(err, w)
		return
	}

	req.Host = strings.Split(r.RemoteAddr, ":")[0]
	res, err := rr.srv.Register(req)
	if err != nil {
		handleErr(err, w)
		return
	}

	// waiting for all workers to announce themselves
	// if amount of registered workers is insufficient - the request is rejected
	registered := <-rr.regChan

	resByte, err := json.Marshal(res)
	if err != nil {
		handleErr(err, w)
		return
	}

	if registered {
		w.WriteHeader(200)
		w.Write(resByte)
	} else {
		w.WriteHeader(500)
		w.Write([]byte("couldn't collect enough workers in time"))
	}
}

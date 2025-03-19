package worker

import (
	"encoding/json"
	"io"
	"net/http"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/worker/internal/domain/worker"
)

type workerRoutes struct {
	wrk worker.Worker
}

func (wr *workerRoutes) handleMap(w http.ResponseWriter, r *http.Request) {
	mapped, err := wr.wrk.Map(r.Body)
	if err != nil {
		handleErr(err, w)
		return
	}

	raw, err := io.ReadAll(mapped)
	if err != nil {
		handleErr(err, w)
		return
	}

	w.WriteHeader(200)
	w.Write(raw)
}

func (wr *workerRoutes) handleReduce(w http.ResponseWriter, r *http.Request) {
	reduced, err := wr.wrk.Reduce(r.Body)
	if err != nil {
		handleErr(err, w)
		return
	}

	raw, err := io.ReadAll(reduced)
	if err != nil {
		handleErr(err, w)
		return
	}

	w.WriteHeader(200)
	w.Write(raw)
}

func (wr *workerRoutes) handleTerminate(w http.ResponseWriter, r *http.Request) {
	var (
		msg mr.TerminateRequest
	)

	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		handleErr(err, w)
		return
	}

	wr.wrk.Terminate(msg.Message)

	w.WriteHeader(http.StatusNoContent)
}

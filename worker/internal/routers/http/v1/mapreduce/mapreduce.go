package mapreduce

import (
	"io"
	"net/http"

	"github.com/cutlery47/map-reduce/worker/internal/domain"
)

type workerRoutes struct {
	wrk domain.Worker
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
	}

	w.WriteHeader(200)
	w.Write(raw)
}

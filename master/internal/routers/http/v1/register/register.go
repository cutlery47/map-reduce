package register

import (
	"encoding/json"
	"net/http"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/mapreduce/requests"
	"github.com/cutlery47/map-reduce/master/internal/domain/master"
)

type registerRoutes struct {
	mst  *master.Master
	conf mr.Config
}

func (rr *registerRoutes) register(w http.ResponseWriter, r *http.Request) {
	var req requests.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleErr(err, w)
		return
	}

	// this was here for some reason:
	// req.Host = strings.Split(r.RemoteAddr, ":")[0]

	// registration handling
	role, err := rr.mst.Register(req)
	if err != nil {
		handleErr(err, w)
		return
	}

	json.NewEncoder(w).Encode(requests.RegisterResponse{
		Role: *role,
	})
	w.WriteHeader(http.StatusOK)
}

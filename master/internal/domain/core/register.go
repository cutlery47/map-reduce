package core

import (
	"errors"
	"fmt"
	"sync"
	"time"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	log "github.com/sirupsen/logrus"
)

// worker registration handler
type RegisterHandler struct {
	mapAddrs, redAddrs []mr.Addr // slices for storing registered worker addresses
	mu                 sync.Mutex
	conf               mr.Config
}

func NewRegisterHandler(conf mr.Config) *RegisterHandler {
	return &RegisterHandler{
		mapAddrs: []mr.Addr{},
		redAddrs: []mr.Addr{},
		mu:       sync.Mutex{},
		conf:     conf,
	}
}

// handles incoming registration requests
// determines whether worker is either a mapper or a reducer
func (rh *RegisterHandler) Handle(req mr.RegisterRequest) (*mr.Role, error) {
	rh.mu.Lock()
	defer rh.mu.Unlock()

	var (
		mappers  = len(rh.mapAddrs)
		reducers = len(rh.redAddrs)
	)

	switch {
	case mappers < rh.conf.Mappers:
		// assign worker to mapper, if possible
		rh.mapAddrs = append(rh.mapAddrs, mr.Addr{
			Host: req.Addr.Host,
			Port: req.Addr.Port,
		})
		return &mr.Mapper, nil
	case reducers < rh.conf.Reducers:
		// assign worker to reducer, if possible
		rh.mapAddrs = append(rh.redAddrs, mr.Addr{
			Host: req.Addr.Host,
			Port: req.Addr.Port,
		})
		return &mr.Reducer, nil
	}

	return nil, errors.New("TI XYESOS")
}

// waits for workers to connect
// returns error if not enough workers have connected
func (rh *RegisterHandler) Collect() ([]mr.Addr, []mr.Addr, error) {
	var (
		curr  = 0
		total = rh.conf.Mappers + rh.conf.Reducers // total amount of possible workers
		timr  = time.NewTimer(rh.conf.RegisterDur)
	)

	for curr != total {
		// count current registered workers
		rh.mu.Lock()
		curr = len(rh.mapAddrs) + len(rh.redAddrs)
		rh.mu.Unlock()

		select {
		case <-timr.C:
			// timed out
			return nil, nil, fmt.Errorf("expected %v workers, connected: %v", total, curr)
		default:
			log.Infoln("currently connected:", curr)
		}

		time.Sleep(1 * time.Second)
	}

	return rh.mapAddrs, rh.redAddrs, nil
}

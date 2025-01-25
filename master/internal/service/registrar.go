package service

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cutlery47/map-reduce/mapreduce"
)

const (
	isMapper  = 0
	isReducer = 1
)

type registrar struct {
	// atomic worker connection count
	mapCnt atomic.Int64
	redCnt atomic.Int64

	mapAddrs *[]addr
	redAddrs *[]addr

	// mutexes for restricting concurrent access to slices
	mapMu *sync.Mutex
	redMu *sync.Mutex

	conf mapreduce.MasterRegistrarConfig
}

func newRegistrar(conf mapreduce.MasterRegistrarConfig, mapAddrs, redAddrs *[]addr) *registrar {
	return &registrar{
		mapCnt:   atomic.Int64{},
		redCnt:   atomic.Int64{},
		mapAddrs: mapAddrs,
		redAddrs: redAddrs,
		mapMu:    &sync.Mutex{},
		redMu:    &sync.Mutex{},
		conf:     conf,
	}
}

// Registering incoming workers
func (reg *registrar) register(req mapreduce.WorkerRegisterRequest) (int, error) {
	// load amounts of currectly connected mappers and reducers
	mappers := reg.mapCnt.Load()
	reducers := reg.redCnt.Load()

	// assign worker to mapper if possible
	if mappers < int64(reg.conf.Mappers) {
		// check if any concurrent goroutine has already updated mapper count
		for mappers < int64(reg.conf.Mappers) && !reg.mapCnt.CompareAndSwap(mappers, mappers+1) {
			// try again after timeout
			time.Sleep(time.Millisecond)
			mappers = reg.mapCnt.Load()
		}

		// add mapper address
		reg.mapMu.Lock()
		*reg.mapAddrs = append(*reg.mapAddrs, addr{Port: req.Port, Host: req.Host})
		reg.mapMu.Unlock()

		return isMapper, nil
	}

	// assign worker to reducer if possible
	if reducers < int64(reg.conf.Reducers) {
		// check if any concurrent goroutine has already updated reducer count
		for reducers < int64(reg.conf.Reducers) && !reg.redCnt.CompareAndSwap(reducers, reducers+1) {
			// try again
			time.Sleep(time.Millisecond)
			reducers = reg.redCnt.Load()
		}

		// add reducer address
		reg.redMu.Lock()
		*reg.redAddrs = append(*reg.redAddrs, addr{Port: req.Port, Host: req.Host})
		reg.redMu.Unlock()

		return isReducer, nil
	}

	return -1, errors.New("received registration request from extra worker")
}

func (reg *registrar) collectWorkers(mappers, reducers int) (int, error) {
	total := mappers + reducers
	deadline := time.Now().Add(reg.conf.RegisterDuration)

	for cnt := reg.mapCnt.Load() + reg.redCnt.Load(); cnt != int64(total); cnt = reg.mapCnt.Load() + reg.redCnt.Load() {
		if time.Now().After(deadline) {
			return int(cnt), fmt.Errorf("expected %v workers, connected: %v", total, cnt)
		}

		log.Println("connected:", cnt)
		time.Sleep(reg.conf.CollectInterval)
	}

	log.Println("all workers connected")
	log.Println("mappers:", reg.mapAddrs)
	log.Println("reducers:", reg.redAddrs)

	return total, nil
}

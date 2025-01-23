package service

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cutlery47/map-reduce/mapreduce"
)

type registrar struct {
	// atomic worker connection count
	cnt atomic.Int64

	mapAddrs *[]addr
	redAddrs *[]addr

	// mutexes for restricting concurrent access to slices
	mapMu *sync.Mutex
	redMu *sync.Mutex

	conf mapreduce.RegistrarConfig
}

func newRegistrar(conf mapreduce.RegistrarConfig, mapAddrs, redAddrs *[]addr) *registrar {
	return &registrar{
		cnt:      atomic.Int64{},
		mapAddrs: mapAddrs,
		redAddrs: redAddrs,
		mapMu:    &sync.Mutex{},
		redMu:    &sync.Mutex{},
		conf:     conf,
	}
}

// Registering incoming workers
func (reg *registrar) register(req mapreduce.WorkerRegisterRequest) error {
	// assign a role for each worker (round robin)
	// each odd worker - is a mapper
	// each even worker - is a reducer
	cur := reg.cnt.Add(1)
	if cur%2 == 0 {
		reg.mapMu.Lock()
		*reg.mapAddrs = append(*reg.mapAddrs, addr{Port: req.Port, Host: req.Host})
		reg.mapMu.Unlock()
	} else {
		reg.redMu.Lock()
		*reg.redAddrs = append(*reg.redAddrs, addr{Port: req.Port, Host: req.Host})
		reg.redMu.Unlock()
	}

	return nil
}

func (reg *registrar) collectWorkers(mappers, reducers int) (int, error) {
	total := mappers + reducers
	deadline := time.Now().Add(reg.conf.RegisterDuration)

	for cnt := reg.cnt.Load(); cnt != int64(total); cnt = reg.cnt.Load() {
		if time.Now().After(deadline) {
			return int(cnt), fmt.Errorf("expected %v workers, connected: %v", total, cnt)
		}

		log.Println("connected:", cnt)
		time.Sleep(reg.conf.CollectInterval)
	}

	log.Println("all workers connected")

	return total, nil
}

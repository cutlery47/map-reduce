package service

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cutlery47/map-reduce/mapreduce"
)

type Service interface {
	Register(req mapreduce.WorkerRegisterRequest) error
}

type MasterService struct {
	// http client for contacting workers
	cl *http.Client

	// slice of assigned mappers
	mapperAddrs []addr
	// slice of assigned reducers
	reducerAddrs []addr

	// atomic worker connection count
	cnt atomic.Int64
	// mutexes for restricting concurrent access to slices
	mapMu *sync.Mutex
	redMu *sync.Mutex

	// channel for signaling that worker collection should start
	startChan chan struct{}

	// Once object for initiating registration once in concurrent environment
	once sync.Once

	// config
	conf mapreduce.Config
}

// MasterService constructor
func NewMasterService(conf mapreduce.Config, cl *http.Client) *MasterService {
	ms := &MasterService{
		cl:           cl,
		cnt:          atomic.Int64{},
		mapperAddrs:  []addr{},
		reducerAddrs: []addr{},
		mapMu:        &sync.Mutex{},
		redMu:        &sync.Mutex{},
		conf:         conf,
		once:         sync.Once{},
		startChan:    make(chan struct{}),
	}

	return ms
}

// Registering incoming workers
func (ms *MasterService) Register(req mapreduce.WorkerRegisterRequest) error {
	ms.once.Do(func() {
		ms.startChan <- struct{}{}
	})

	// assign a role for each worker (round robin)
	// each odd worker - is a mapper
	// each even worker - is a reducer
	cur := ms.cnt.Add(1)
	if cur%2 == 0 {
		ms.mapMu.Lock()
		ms.mapperAddrs = append(ms.mapperAddrs, addr{Port: req.Port, Host: req.Host})
		ms.mapMu.Unlock()
	} else {
		ms.redMu.Lock()
		ms.reducerAddrs = append(ms.reducerAddrs, addr{Port: req.Port, Host: req.Host})
		ms.redMu.Unlock()
	}

	return nil
}

// Handle registered workers
func (ms *MasterService) HandleWorkers(errChan chan<- error, regChan chan<- bool) {
	// waiting for first request to hit
	<-ms.startChan

	total := ms.conf.Mappers + ms.conf.Reducers
	deadline := time.Now().Add(ms.conf.RegisterDuration)

	for cnt := ms.cnt.Load(); cnt != int64(total); cnt = ms.cnt.Load() {
		if time.Now().After(deadline) {
			for range cnt {
				regChan <- false
			}
			errChan <- fmt.Errorf("expected %v workers, connected: %v", total, cnt)
			return
		}

		log.Println("connected:", cnt)
		time.Sleep(ms.conf.CollectTimeout)
	}
	log.Println("all workers connected")

	// sending confirmation responses for all pending registration requests
	for range total {
		regChan <- true
	}

	// split current file into parts = amount of workers
	files, err := splitFile(ms.conf.FileLocation, ms.conf.ChunkLocation, ms.conf.ResultLocation, ms.conf.Mappers)
	if err != nil {
		errChan <- fmt.Errorf("SplitFile: %v", err)
		return
	}

	var mapResults [][]byte
	var redResults [][]byte

	// Senging tasks to mappers
	for i, addr := range ms.mapperAddrs {
		res, err := ms.cl.Post(fmt.Sprintf("http://%v:%v/map", addr.Host, addr.Port), "text/plain", files[i])
		if err != nil {
			errChan <- err
			return
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			errChan <- err
			return
		}

		// storing map results
		mapResults = append(mapResults, body)
	}

	// Sending tasks to reducers
	for i, addr := range ms.reducerAddrs {
		res, err := ms.cl.Post(fmt.Sprintf("http://%v:%v/reduce", addr.Host, addr.Port), "application/json", bytes.NewReader(mapResults[i]))
		if err != nil {
			errChan <- err
			return
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			errChan <- err
			return
		}

		// storing reduce results
		redResults = append(redResults, body)
	}

	for i, res := range redResults {
		fileName := fmt.Sprintf("res_%v", i)
		if err := execCreateAndWriteFile(fileName, ms.conf.ResultLocation, res); err != nil {
			errChan <- err
			return
		}
	}

	time.Sleep(1 * time.Hour)
	// shutdown signal
	errChan <- err
}

type addr struct {
	Port string
	Host string
}

package service

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/cutlery47/map-reduce/mapreduce"
)

type Service interface {
	Register(req mapreduce.WorkerRegisterRequest) (mapreduce.WorkerRegisterResponse, error)
	HandleWorkers(errChan chan<- error, regChan chan<- bool)
}

type MasterService struct {
	reg *registrar
	tp  taskProducer

	// config
	conf mapreduce.MasterConfig

	// Once object for initiating registration once in concurrent environment
	once sync.Once
	// channel for signaling that worker collection should start
	startChan chan struct{}
}

func NewMasterService(conf mapreduce.MasterConfig, cl *http.Client) (*MasterService, error) {
	mapAddrs := []addr{}
	redAddrs := []addr{}

	reg := newRegistrar(conf.MasterRegistrarConfig, &mapAddrs, &redAddrs)

	var tp taskProducer

	switch conf.ProducerType {
	case "HTTP":
		tp = NewHTTPTaskProducer(cl, &mapAddrs, &redAddrs)
	case "QUEUE":
		prod, err := NewRabbitTaskProducer(conf.RabbitConfig)
		if err != nil {
			return nil, fmt.Errorf("NewRabbitTaskProducer: %v", err)
		}
		tp = prod
	default:
		return nil, errors.New("undefined task producer")
	}

	return &MasterService{
		reg:       reg,
		tp:        tp,
		conf:      conf,
		once:      sync.Once{},
		startChan: make(chan struct{}),
	}, nil
}

func (ms *MasterService) Register(req mapreduce.WorkerRegisterRequest) (mapreduce.WorkerRegisterResponse, error) {
	ms.once.Do(func() {
		ms.startChan <- struct{}{}
	})

	res := mapreduce.WorkerRegisterResponse{}

	wType, err := ms.reg.register(req)
	if err != nil {
		return res, nil
	}

	switch wType {
	case isMapper:
		res.Type = mapreduce.MapperType
		return res, nil
	case isReducer:
		res.Type = mapreduce.ReducerType
		return res, nil
	}

	panic("idk")
}

// Handle registered workers
func (ms *MasterService) HandleWorkers(errChan chan<- error, regChan chan<- bool) {
	// waiting for first request to hit
	<-ms.startChan

	cnt, err := ms.reg.collectWorkers(ms.conf.Mappers, ms.conf.Reducers)
	if err != nil {
		for range cnt {
			regChan <- false
		}
		errChan <- fmt.Errorf("ms.reg.collectWorkers: %v", err)
		return
	}

	// sending confirmation responses for all pending registration requests
	for range cnt {
		regChan <- true
	}

	var b strings.Builder

	fileLocation := createNestedDirString(b, ms.conf.FileDirectory, ms.conf.FileName)
	chunkLocation := createNestedDirString(b, ms.conf.FileDirectory, ms.conf.ChunkDirectory)
	resultLocation := createNestedDirString(b, ms.conf.FileDirectory, ms.conf.ResultDirectory)

	// split current file into parts = amount of workers
	files, err := splitFile(fileLocation, chunkLocation, resultLocation, ms.conf.Mappers)
	if err != nil {
		errChan <- fmt.Errorf("SplitFile: %v", err)
		return
	}

	// Senging tasks to mappers
	mapResults, err := ms.tp.produceMapperTasks(files)
	if err != nil {
		errChan <- fmt.Errorf("ms.sendMapperTasks: %v", err)
		return
	}

	// Sending tasks to reducers
	redResults, err := ms.tp.produceReducerTasks(mapResults)
	if err != nil {
		errChan <- fmt.Errorf("ms.sendReducerTasks: %v", err)
		return
	}

	// Handling reducer results
	for i, res := range redResults {
		fileName := fmt.Sprintf("res_%v", i)
		if err := execCreateAndWriteFile(fileName, resultLocation, res); err != nil {
			errChan <- err
			return
		}
	}

	// shutdown signal
	errChan <- fmt.Errorf("finished map-reduce")
}

type addr struct {
	Port string
	Host string
}

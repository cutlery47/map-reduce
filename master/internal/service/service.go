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
	Register(req mapreduce.WorkerRegisterRequest) error
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

	reg := newRegistrar(conf.RegistrarConfig, &mapAddrs, &redAddrs)

	var tp taskProducer

	switch conf.ProducerType {
	case "HTTP":
		tp = NewHTTPTaskProducer(cl, &mapAddrs, &redAddrs)
	case "RABBIT":
		prod, err := NewRabbitTaskProducer(conf.ProducerConfig)
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

func (ms *MasterService) Register(req mapreduce.WorkerRegisterRequest) error {
	ms.once.Do(func() {
		ms.startChan <- struct{}{}
	})

	return ms.reg.register(req)
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
	}

	// Sending tasks to reducers
	redResults, err := ms.tp.produceReducerTasks(mapResults)
	if err != nil {
		errChan <- fmt.Errorf("ms.sendReducerTasks: %v", err)
	}

	// Handling reducer results
	for i, res := range redResults {
		fileName := fmt.Sprintf("res_%v", i)
		if err := execCreateAndWriteFile(fileName, resultLocation, res); err != nil {
			errChan <- err
			return
		}
	}

	fmt.Println("finished map-reduce")
	// shutdown signal
	errChan <- err
}

type addr struct {
	Port string
	Host string
}

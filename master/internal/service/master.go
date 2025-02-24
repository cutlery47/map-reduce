package service

import (
	"fmt"
	"strings"
	"sync"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/master/internal/core"
	prod "github.com/cutlery47/map-reduce/master/internal/service/producer"
	"github.com/cutlery47/map-reduce/master/internal/util"
)

type Master struct {
	reg *core.Registrar
	tp  prod.TaskProducer

	// config
	conf mr.MstSvcConf

	// Once object for initiating registration once in concurrent environment
	once sync.Once
	// channel for signaling that worker collection should start
	startChan chan struct{}
}

func New(conf mr.MstSvcConf, reg *core.Registrar, tp prod.TaskProducer) (*Master, error) {
	return &Master{
		reg:       reg,
		tp:        tp,
		conf:      conf,
		once:      sync.Once{},
		startChan: make(chan struct{}),
	}, nil
}

func (ms *Master) Register(req mr.WrkRegReq) (mr.WrkRegRes, error) {
	ms.once.Do(func() {
		ms.startChan <- struct{}{}
	})

	res := mr.WrkRegRes{}

	wType, err := ms.reg.Register(req)
	if err != nil {
		return res, nil
	}

	switch wType {
	case core.IsMapper:
		res.Type = mr.MapperType
		return res, nil
	case core.IsReducer:
		res.Type = mr.ReducerType
		return res, nil
	}

	panic("idk")
}

// Handle registered workers
func (ms *Master) HandleWorkers(regChan chan<- bool) error {
	// waiting for first request to hit
	<-ms.startChan

	cnt, err := ms.reg.CollectWorkers(ms.conf.Mappers, ms.conf.Reducers)
	if err != nil {
		for range cnt {
			regChan <- false
		}
		return fmt.Errorf("ms.reg.collectWorkers: %v", err)
	}

	// sending confirmation responses for all pending registration requests
	for range cnt {
		regChan <- true
	}

	var b strings.Builder

	fileLocation := util.CreateNestedDirString(b, ms.conf.FileDirectory, ms.conf.FileName)
	chunkLocation := util.CreateNestedDirString(b, ms.conf.FileDirectory, ms.conf.ChunkDirectory)
	resultLocation := util.CreateNestedDirString(b, ms.conf.FileDirectory, ms.conf.ResultDirectory)

	// split current file into parts = amount of workers
	files, err := util.SplitFile(fileLocation, chunkLocation, resultLocation, ms.conf.Mappers)
	if err != nil {
		return fmt.Errorf("SplitFile: %v", err)
	}

	// Senging tasks to mappers
	mapResults, err := ms.tp.ProduceMapperTasks(files)
	if err != nil {
		return fmt.Errorf("ms.sendMapperTasks: %v", err)
	}

	// Sending tasks to reducers
	redResults, err := ms.tp.ProduceReducerTasks(mapResults)
	if err != nil {
		return fmt.Errorf("ms.sendReducerTasks: %v", err)
	}

	// Handling reducer results
	for i, res := range redResults {
		fileName := fmt.Sprintf("res_%v", i)
		if err := util.ExecCreateAndWriteFile(fileName, resultLocation, res); err != nil {
			return err
		}
	}

	// shutdown signal
	return fmt.Errorf("finished map-reduce")
}

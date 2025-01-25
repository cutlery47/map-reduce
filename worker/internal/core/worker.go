package core

import (
	"fmt"
	"io"
	"log"

	"github.com/cutlery47/map-reduce/mapreduce"
)

var mapFunc, reduceFunc = mapreduce.MapReduce()

type Worker interface {
	Map(reader io.Reader) (any, error)
	Reduce(mapped any) (any, error)
}

type DefaultWorker struct {
	// channel for signaling that the worker has finished execution
	endChan chan<- struct{}

	// config
	conf mapreduce.WorkerConfig
}

func NewDefaultWorker(conf mapreduce.WorkerConfig, endChan chan<- struct{}) *DefaultWorker {
	return &DefaultWorker{
		conf:    conf,
		endChan: endChan,
	}
}

func (dw *DefaultWorker) Map(reader io.Reader) (any, error) {
	log.Println("in map")

	mapResult, err := mapFunc(reader)
	if err != nil {
		return nil, err
	}

	log.Println("finished map")
	dw.endChan <- struct{}{}

	return mapResult, nil
}

func (dw *DefaultWorker) Reduce(result interface{}) (any, error) {
	log.Println("in reduce")

	mapResult, ok := result.(mapreduce.MyMapResult)
	if !ok {
		return nil, fmt.Errorf("map result incompatible with reduce func")
	}

	redResult, err := reduceFunc(mapResult)
	if err != nil {
		return nil, err
	}

	log.Println("finished reduce")
	dw.endChan <- struct{}{}

	return redResult, nil
}

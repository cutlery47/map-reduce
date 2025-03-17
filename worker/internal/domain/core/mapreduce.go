package core

import (
	"fmt"
	"io"
	"log"

	mr "github.com/cutlery47/map-reduce/mapreduce"
)

var mapFunc, reduceFunc = mr.MapReduce()

type MapReduceHandler struct {
	conf mr.Config
}

func NewMapReduceHandler(conf mr.Config) (*MapReduceHandler, error) {
	return &MapReduceHandler{
		conf: conf,
	}, nil
}

func (mrh *MapReduceHandler) Map(reader io.Reader) (any, error) {
	log.Println("in map")

	mapResult, err := mapFunc(reader)
	if err != nil {
		return nil, err
	}

	log.Println("finished map")

	return mapResult, nil
}

func (mrh *MapReduceHandler) Reduce(result interface{}) (any, error) {
	log.Println("in reduce")

	mapResult, ok := result.(mr.MyMapResult)
	if !ok {
		return nil, fmt.Errorf("map result incompatible with reduce func")
	}

	redResult, err := reduceFunc(mapResult)
	if err != nil {
		return nil, err
	}

	log.Println("finished reduce")

	return redResult, nil
}

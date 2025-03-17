package domain

import (
	"errors"
	"fmt"
	"sync"
	"time"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/master/internal/domain/core"
	prod "github.com/cutlery47/map-reduce/master/internal/domain/producer"
)

type Master struct {
	tp prod.TaskProducer
	rh *core.RegisterHandler
	fh *core.FileHandler

	once    sync.Once     // once object for initiating worker registration once
	startCh chan struct{} // channel for signaling that worker collection should start

	conf mr.Config
}

func NewMaster(conf mr.Config, tp prod.TaskProducer) *Master {
	var (
		startCh = make(chan struct{})
	)

	return &Master{
		tp:      tp,
		rh:      core.NewRegisterHandler(conf),
		fh:      core.NewFileHandler(conf),
		conf:    conf,
		once:    sync.Once{},
		startCh: startCh,
	}
}

// primary master logic
func (ms *Master) Work() error {
	// preparing directories to store files in
	if err := ms.fh.CreateDirs(); err != nil {
		return err
	}

	var timr = time.NewTimer(ms.conf.MasterRequestAwaitDur)
	// waiting for first request to hit
	select {
	case <-ms.startCh:
		// start registering workers
	case <-timr.C:
		// no requests received
		return errors.New("no requests received")
	}

	// waiting for all workers to connect
	mappers, reducers, err := ms.rh.Collect()
	if err != nil {
		return err
	}

	// split current file into chunks = amount of mappers
	chunks, err := ms.fh.Split(ms.conf.Mappers)
	if err != nil {
		return fmt.Errorf("SplitFile: %v", err)
	}

	// senging tasks to mappers
	mapResults, err := ms.tp.ProduceMapperTasks(chunks, mappers)
	if err != nil {
		return fmt.Errorf("ms.sendMapperTasks: %v", err)
	}

	// sending tasks to reducers
	redResults, err := ms.tp.ProduceReducerTasks(mapResults, reducers)
	if err != nil {
		return fmt.Errorf("ms.sendReducerTasks: %v", err)
	}

	// handling reducer results
	for i, res := range redResults {
		if err := ms.fh.CreateResult(fmt.Sprintf("res_%v", i), res); err != nil {
			return err
		}
	}

	// shutdown signal
	return fmt.Errorf("finished map-reduce")
}

// passes incoming registration requests to registrar
// also, initializes collection process on first request
func (ms *Master) Register(req mr.RegisterRequest) (*mr.Role, error) {
	ms.once.Do(func() {
		ms.startCh <- struct{}{}
	})

	return ms.rh.Handle(req)
}

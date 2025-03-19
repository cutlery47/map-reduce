package master

import (
	"errors"
	"fmt"
	"sync"
	"time"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/master/internal/domain/core"
	prod "github.com/cutlery47/map-reduce/master/internal/domain/producer"
	log "github.com/sirupsen/logrus"
)

// default master implementation
type Master struct {
	tp prod.TaskProducer
	rh *core.RegisterHandler
	fh *core.FileHandler

	once    sync.Once     // once object for initiating worker registration once
	startCh chan struct{} // channel for signaling that worker collection should start

	conf mr.Config
}

func New(conf mr.Config, tp prod.TaskProducer) (*Master, error) {
	return &Master{
		tp:      tp,
		rh:      core.NewRegisterHandler(conf),
		fh:      core.NewFileHandler(conf),
		conf:    conf,
		once:    sync.Once{},
		startCh: make(chan struct{}),
	}, nil
}

// primary master logic
func (sm *Master) Run() error {
	// preparing directories to store files in
	log.Infoln("[MASTER] creating necessary directories...")
	if err := sm.fh.CreateDirs(); err != nil {
		return err
	}

	var timr = time.NewTimer(sm.conf.MasterRequestAwaitDur)
	log.Infoln("[MASTER] waiting for worker to register...")

	// waiting for first request to hit
	select {
	case <-sm.startCh:
		log.Infoln("[MASTER] registration initiated...")
	case <-timr.C:
		// no requests received
		return errors.New("no requests received")
	}

	// waiting for all workers to connect
	mappers, reducers, err := sm.rh.Collect()
	if err != nil {
		return err
	}

	// split current file into chunks = amount of mappers
	chunks, err := sm.fh.Split(sm.conf.Mappers)
	if err != nil {
		return fmt.Errorf("SplitFile: %v", err)
	}

	// senging tasks to mappers
	mapResults, err := sm.tp.ProduceMapperTasks(chunks, mappers)
	if err != nil {
		return fmt.Errorf("sm.sendMapperTasks: %v", err)
	}

	// sending tasks to reducers
	redResults, err := sm.tp.ProduceReducerTasks(mapResults, reducers)
	if err != nil {
		return fmt.Errorf("sm.sendReducerTasks: %v", err)
	}

	// handling reducer results
	for i, res := range redResults {
		if err := sm.fh.CreateResult(fmt.Sprintf("res_%v", i), res); err != nil {
			return err
		}
	}

	// shutdown signal
	return fmt.Errorf("finished map-reduce")
}

// passes incoming registration requests to registrar
// also, initializes collection process on first request
func (sm *Master) Register(req mr.RegisterRequest) (*mr.Role, error) {
	sm.once.Do(func() {
		sm.startCh <- struct{}{}
	})

	return sm.rh.Register(req)
}

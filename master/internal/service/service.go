package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cutlery47/map-reduce/mapreduce"
)

type Service interface {
	Register(req mapreduce.WorkerRegisterRequest) error
}

type WorkerService struct {
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

	// channel for sending error to the httpserver
	errChan chan<- error
	// channel for signaling that worker collection should start
	startChan chan struct{}

	once sync.Once

	conf mapreduce.Config
}

// WorkerService constructor
func NewWorkerService(conf mapreduce.Config, cl *http.Client, errChan chan<- error) *WorkerService {
	ws := &WorkerService{
		cl:           cl,
		cnt:          atomic.Int64{},
		mapperAddrs:  []addr{},
		reducerAddrs: []addr{},
		mapMu:        &sync.Mutex{},
		redMu:        &sync.Mutex{},
		conf:         conf,
		once:         sync.Once{},
		errChan:      errChan,
		startChan:    make(chan struct{}),
	}

	go ws.HandleWorkers()
	return ws
}

// Registering incoming workers
func (ws *WorkerService) Register(req mapreduce.WorkerRegisterRequest) error {
	ws.once.Do(func() {
		ws.startChan <- struct{}{}
	})

	cur := ws.cnt.Add(1)
	if cur%2 == 0 {
		ws.mapMu.Lock()
		ws.mapperAddrs = append(ws.mapperAddrs, addr{Port: req.Port, Host: req.Host})
		ws.mapMu.Unlock()
	} else {
		ws.redMu.Lock()
		ws.reducerAddrs = append(ws.reducerAddrs, addr{Port: req.Port, Host: req.Host})
		ws.redMu.Unlock()
	}

	return nil
}

// Handle registered workers
func (ws *WorkerService) HandleWorkers() {
	// waiting for first request to hit
	<-ws.startChan

	total := ws.conf.Mappers + ws.conf.Reducers
	timer := time.NewTimer(ws.conf.RegisterTimeout)
	for cnt := ws.cnt.Load(); cnt != int64(total); cnt = ws.cnt.Load() {
		select {
		case <-timer.C:
			ws.errChan <- fmt.Errorf("expected %v workers, connected: %v", total, cnt)
			return
		default:
			log.Println("connected:", cnt)
			time.Sleep(ws.conf.CollectTimeout)
		}
	}
	log.Println("all workers connected")

	// split current file into parts = amount of workers
	files, err := ws.SplitFile("file.txt", ws.conf.Mappers)
	if err != nil {
		ws.errChan <- fmt.Errorf("SplitFile: %v", err)
		return
	}

	var mapResults [][]byte
	var redResults [][]byte

	// Senging tasks to mappers
	for i, addr := range ws.mapperAddrs {
		res, err := ws.cl.Post(fmt.Sprintf("http://%v:%v/map", addr.Host, addr.Port), "text/plain", files[i])
		if err != nil {
			ws.errChan <- err
			return
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			ws.errChan <- err
			return
		}

		// storing map results
		mapResults = append(mapResults, body)
	}

	// Sending tasks to reducers
	for i, addr := range ws.reducerAddrs {
		res, err := ws.cl.Post(fmt.Sprintf("http://%v:%v/reduce", addr.Host, addr.Port), "application/json", bytes.NewReader(mapResults[i]))
		if err != nil {
			ws.errChan <- err
			return
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			ws.errChan <- err
			return
		}

		// storing reduce results
		redResults = append(redResults, body)
	}

	// return and shutdown
	err = ws.ReturnResults(redResults)
	ws.errChan <- err
}

// parses reduce results and outpts them as a map
func (ws *WorkerService) ReturnResults(redResults [][]byte) error {
	hashMap := make(map[string]int)

	for _, res := range redResults {
		sliceRes := []string{}
		json.Unmarshal(res, &sliceRes)
		fmt.Println(sliceRes)

		for _, el := range sliceRes {
			elSplit := strings.Split(el, ":")
			key := elSplit[0]
			val, err := strconv.Atoi(elSplit[1])
			if err != nil {
				return err
			}

			if v, ok := hashMap[key]; !ok {
				hashMap[key] = val
			} else {
				hashMap[key] = v + val
			}
		}
	}

	return nil
}

func (ws *WorkerService) SplitFile(filename string, parts int) ([]io.Reader, error) {
	readers := []io.Reader{}

	err := os.Mkdir("chunks", 0777)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			return nil, fmt.Errorf("os.Mkdir: %v", err)
		}
	}

	bash := "split"
	arg0, arg1 := "-n", strconv.Itoa(parts)
	arg2 := "-d"
	arg3, arg4 := "file.txt", "chunks/"

	cmd := exec.Command(bash, arg0, arg1, arg2, arg3, arg4)
	cmd.Stderr = os.Stdin
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	chunks, err := os.ReadDir("chunks")
	if err != nil {
		return nil, fmt.Errorf("os.ReadDir: %v", err)
	}

	for _, chunk := range chunks {
		fd, err := os.Open(fmt.Sprintf("chunks/%v", chunk.Name()))
		if err != nil {
			return nil, err
		}

		readers = append(readers, fd)
	}

	return readers, nil
}

type addr struct {
	Port string
	Host string
}

package httpmaster

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/cutlery47/map-reduce/master/internal/core"
)

type httpTaskProducer struct {
	// http client for contacting workers
	cl *http.Client

	mapAddrs, redAddrs *[]core.Addr
}

// HTTPWorkerHandler constructor
func NewHTTPTaskProducer(cl *http.Client, mapAddrs, redAddrs *[]core.Addr) *httpTaskProducer {
	return &httpTaskProducer{
		mapAddrs: mapAddrs,
		redAddrs: redAddrs,
		cl:       cl,
	}
}

// // produce tasks for mappers over http
func (htp *httpTaskProducer) ProduceMapperTasks(in []io.Reader) ([][]byte, error) {
	return htp.produce(in, *htp.mapAddrs, true)
}

// produce tasks for reducers over http
func (htp *httpTaskProducer) ProduceReducerTasks(in [][]byte) ([][]byte, error) {
	input := make([]io.Reader, len(in))

	for i, res := range in {
		input[i] = bytes.NewReader(res)
	}

	return htp.produce(input, *htp.redAddrs, false)
}

func (htp *httpTaskProducer) produce(input []io.Reader, addrs []core.Addr, toMapper bool) ([][]byte, error) {
	output := [][]byte{}

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	// channel for passing errors from goroutines
	errChan := make(chan error, len(*htp.mapAddrs))

	// check where task is bound to (toMapper value)
	// if toMapper == map -- send task to the map endpoint of a worker
	// else -- to the reduce endpoint
	var suffix string
	if toMapper == true {
		suffix = "map"
	} else {
		suffix = "reduce"
	}

	// asynchronously sending each worker a task
	for i, addr := range addrs {
		wg.Add(1)

		go func() {
			defer wg.Done()

			res, err := htp.cl.Post(fmt.Sprintf("http://%v:%v/%v", addr.Host, addr.Port, suffix), "application/json", input[i])
			if err != nil {
				errChan <- fmt.Errorf("ms.cl.Post: %v", err)
				return
			}

			if res == nil {
				errChan <- fmt.Errorf("%v:%v has disconnected during registration", addr.Host, addr.Port)
				return
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				errChan <- fmt.Errorf("ms.cl.Post: %v", err)
				return
			}

			// storing reduce results
			mu.Lock()
			output = append(output, body)
			mu.Unlock()
		}()
	}

	wg.Wait()

	// check if any errors occurred when sending tasks
	select {
	case err := <-errChan:
		log.Println("received error:", err)
		return nil, err
	default:
		log.Printf("all %v tasks sent successfully\n", suffix)
	}

	return output, nil
}

package producer

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	mr "github.com/cutlery47/map-reduce/mapreduce"
)

type httpTaskProducer struct {
	// http client for contacting workers
	cl *http.Client
}

func New(conf mr.Config) *httpTaskProducer {
	return &httpTaskProducer{
		cl: http.DefaultClient,
	}
}

// produce tasks for mappers over http
func (htp *httpTaskProducer) ProduceMapperTasks(input []io.Reader, mappers []mr.Addr) ([][]byte, error) {
	return htp.produce(input, mappers, true)
}

// produce tasks for reducers over http
func (htp *httpTaskProducer) ProduceReducerTasks(in [][]byte, reducers []mr.Addr) ([][]byte, error) {
	input := make([]io.Reader, len(in))

	for i, res := range in {
		input[i] = bytes.NewReader(res)
	}

	return htp.produce(input, reducers, false)
}

func (htp *httpTaskProducer) produce(input []io.Reader, addrs []mr.Addr, toMapper bool) ([][]byte, error) {
	var (
		output [][]byte
		wg     sync.WaitGroup
		mu     sync.Mutex
		suffix string
		errCh  = make(chan error, len(addrs)) // channel for passing errors from goroutines
	)

	// check where task is bound to (toMapper value)
	// if toMapper == map -- send task to the map endpoint of a worker
	// else -- to the reduce endpoint
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
				errCh <- fmt.Errorf("ms.cl.Post: %v", err)
				return
			}

			if res == nil {
				errCh <- fmt.Errorf("%v:%v returned a nil response", addr.Host, addr.Port)
				return
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				errCh <- fmt.Errorf("ms.cl.Post: %v", err)
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
	case err := <-errCh:
		log.Println("received error:", err)
		return nil, err
	default:
		log.Printf("all %v tasks sent successfully\n", suffix)
	}

	return output, nil
}

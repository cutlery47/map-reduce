package producer

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/mapreduce/models"
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
func (htp *httpTaskProducer) ProduceMapperTasks(input []io.Reader, mappers []models.Addr) ([][]byte, error) {
	return htp.produce(input, mappers, true)
}

// produce tasks for reducers over http
func (htp *httpTaskProducer) ProduceReducerTasks(in [][]byte, reducers []models.Addr) ([][]byte, error) {
	input := make([]io.Reader, len(in))

	for i, res := range in {
		input[i] = bytes.NewReader(res)
	}

	return htp.produce(input, reducers, false)
}

func (htp *httpTaskProducer) produce(input []io.Reader, addrs []models.Addr, toMapper bool) ([][]byte, error) {
	var (
		output [][]byte
		wg     sync.WaitGroup
		mu     sync.Mutex
		suffix string
		errCh  = make(chan error, len(addrs)) // channel for passing errors from goroutines
	)

	// check where job is bound to
	if toMapper == true {
		suffix = "map"
	} else {
		suffix = "reduce"
	}

	// function for sending jobs
	sendFunc := func(addr models.Addr, chunk io.Reader) error {
		defer wg.Done()

		var (
			workerAddr = fmt.Sprintf("http://%v:%v/api/v1/worker/%v", addr.Host, addr.Port, suffix)
		)

		res, err := htp.cl.Post(workerAddr, "application/json", chunk)
		if err != nil {
			return err
		}

		if res.StatusCode != http.StatusOK {
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				return err
			}
			return fmt.Errorf("registration failed: %v", string(resBody))
		}

		if res == nil {
			return errors.New("empty worker response")
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		// storing reduce results
		mu.Lock()
		output = append(output, body)
		mu.Unlock()

		return nil
	}

	// sending each worker a job in parallel
	for i, addr := range addrs {
		wg.Add(1)

		go func() {
			err := sendFunc(addr, input[i])
			if err != nil {
				errCh <- err
			}
		}()
	}

	wg.Wait()

	select {
	case err := <-errCh:
		return nil, err
	default:
	}

	return output, nil
}

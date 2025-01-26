package service

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/cutlery47/map-reduce/mapreduce"
	amqp "github.com/rabbitmq/amqp091-go"
)

type taskProducer interface {
	produceMapperTasks(files []io.Reader) ([][]byte, error)
	produceReducerTasks(mapResults [][]byte) ([][]byte, error)
}

type HTTPTaskProducer struct {
	// http client for contacting workers
	cl *http.Client

	mapAddrs, redAddrs *[]addr
}

// HTTPWorkerHandler constructor
func NewHTTPTaskProducer(cl *http.Client, mapAddrs, redAddrs *[]addr) *HTTPTaskProducer {
	return &HTTPTaskProducer{
		mapAddrs: mapAddrs,
		redAddrs: redAddrs,
		cl:       cl,
	}
}

// // produce tasks for mappers over http
func (htp *HTTPTaskProducer) produceMapperTasks(files []io.Reader) ([][]byte, error) {
	return htp.produce(files, *htp.mapAddrs, true)
}

// produce tasks for reducers over http
func (htp *HTTPTaskProducer) produceReducerTasks(mapResults [][]byte) ([][]byte, error) {
	input := make([]io.Reader, len(mapResults))

	for i, res := range mapResults {
		input[i] = bytes.NewReader(res)
	}

	return htp.produce(input, *htp.redAddrs, false)
}

func (htp *HTTPTaskProducer) produce(input []io.Reader, addrs []addr, toMapper bool) ([][]byte, error) {
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

type RabbitTaskProducer struct {
	b *Brocker

	conf mapreduce.RabbitConfig
}

func NewRabbitTaskProducer(conf mapreduce.RabbitConfig) (*RabbitTaskProducer, error) {
	b, err := NewBrocker(conf)
	if err != nil {
		return nil, fmt.Errorf("NewBrocker: %v", err)
	}

	return &RabbitTaskProducer{
		b:    b,
		conf: conf,
	}, nil
}

func (rtp *RabbitTaskProducer) produceMapperTasks(files []io.Reader) ([][]byte, error) {
	q, err := rtp.b.ch.QueueDeclare(
		rtp.conf.MapperQueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	err = rtp.b.ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("123123123"),
		},
	)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ms *RabbitTaskProducer) produceReducerTasks(mapResults [][]byte) ([][]byte, error) {

	return nil, nil
}

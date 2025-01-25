package service

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/cutlery47/map-reduce/mapreduce"
	"github.com/rabbitmq/amqp091-go"
)

type taskProducer interface {
	produceMapperTasks(files []io.Reader) ([][]byte, error)
	produceReducerTasks(mapResults [][]byte) ([][]byte, error)
}

type HTTPTaskProducer struct {
	// http client for contacting workers
	cl *http.Client

	mapAddrs *[]addr
	redAddrs *[]addr
}

// HTTPWorkerHandler constructor
func NewHTTPTaskProducer(cl *http.Client, mapAddrs, redAddrs *[]addr) *HTTPTaskProducer {
	return &HTTPTaskProducer{
		mapAddrs: mapAddrs,
		redAddrs: redAddrs,
		cl:       cl,
	}
}

func (htp *HTTPTaskProducer) produceMapperTasks(files []io.Reader) ([][]byte, error) {
	mapResults := [][]byte{}

	for i, addr := range *htp.mapAddrs {
		res, err := htp.cl.Post(fmt.Sprintf("http://%v:%v/map", addr.Host, addr.Port), "text/plain", files[i])
		if err != nil {
			return nil, fmt.Errorf("ms.cl.Post: %v", err)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("io.ReadAll: %v", err)
		}

		// storing map results
		mapResults = append(mapResults, body)
	}

	return mapResults, nil
}

func (htp *HTTPTaskProducer) produceReducerTasks(mapResults [][]byte) ([][]byte, error) {
	redResults := [][]byte{}

	for i, addr := range *htp.redAddrs {
		res, err := htp.cl.Post(fmt.Sprintf("http://%v:%v/reduce", addr.Host, addr.Port), "application/json", bytes.NewReader(mapResults[i]))
		if err != nil {
			return nil, fmt.Errorf("ms.cl.Post: %v", err)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("ms.cl.Post: %v", err)
		}

		// storing reduce results
		redResults = append(redResults, body)
	}

	return redResults, nil
}

type RabbitTaskProducer struct {
	mapMq, redMq *Brocker
}

func NewRabbitTaskProducer(conf mapreduce.ProducerConfig) (*RabbitTaskProducer, error) {
	mapMqConf, err := mapreduce.NewRabbitConfig(conf.RabbitMapperPath)
	if err != nil {
		return nil, fmt.Errorf("mapreduce.NewRabbitConfig: %v", err)
	}

	redMqConf, err := mapreduce.NewRabbitConfig(conf.RabbitReducerPath)
	if err != nil {
		return nil, fmt.Errorf("mapreduce.NewRabbitConfig: %v", err)
	}

	mapMq, err := NewBrocker(mapMqConf)
	if err != nil {
		return nil, fmt.Errorf("NewBrocker: %v", err)
	}

	redMq, err := NewBrocker(redMqConf)
	if err != nil {
		return nil, fmt.Errorf("NewBrocker: %v", err)
	}

	return &RabbitTaskProducer{
		mapMq: mapMq,
		redMq: redMq,
	}, nil
}

func (rtp *RabbitTaskProducer) produceMapperTasks(files []io.Reader) ([][]byte, error) {
	q, err := rtp.mapMq.ch.QueueDeclare(
		"some",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		buf, _ := io.ReadAll(file)

		rtp.mapMq.ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp091.Publishing{
				ContentType: "text/plain",
				Body:        buf,
			},
		)
	}

	return nil, nil
}

func (ms *RabbitTaskProducer) produceReducerTasks(mapResults [][]byte) ([][]byte, error) {

	return nil, nil
}

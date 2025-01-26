package queueworker

import (
	"fmt"
	"io"

	"github.com/cutlery47/map-reduce/mapreduce"
	amqp "github.com/rabbitmq/amqp091-go"
)

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

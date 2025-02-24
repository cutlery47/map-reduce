package queue

import (
	"io"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitTaskProducer struct {
	ch *amqp.Channel

	conf mr.RabbitConf
}

func NewRabbitTaskProducer(conf mr.RabbitConf, ch *amqp.Channel) (*rabbitTaskProducer, error) {
	return &rabbitTaskProducer{
		ch:   ch,
		conf: conf,
	}, nil
}

func (rtp *rabbitTaskProducer) produceMapperTasks(in []io.Reader) ([][]byte, error) {
	q, err := rtp.ch.QueueDeclare(
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

	err = rtp.ch.Publish(
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

func (ms *rabbitTaskProducer) produceReducerTasks(in [][]byte) ([][]byte, error) {

	return nil, nil
}

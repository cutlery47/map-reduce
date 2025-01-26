package queue

import (
	"fmt"

	"github.com/cutlery47/map-reduce/mapreduce"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Brocker struct {
	ch *amqp.Channel

	conf mapreduce.RabbitConfig
}

func NewBrocker(conf mapreduce.RabbitConfig) (*Brocker, error) {
	url := fmt.Sprintf(
		"amqp://%v:%v@%v:%v/",
		conf.RabbitLogin,
		conf.RabbitPassword,
		conf.RabbitHost,
		conf.RabbitPort,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("amqp.Dial: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("conn.Channel: %v", err)
	}

	return &Brocker{
		ch:   ch,
		conf: conf,
	}, nil
}

func (b *Brocker) DeclareAndConsume(wType string) (<-chan amqp.Delivery, error) {
	var name string

	switch wType {
	case mapreduce.MapperType:
		name = b.conf.MapperQueueName
	case mapreduce.ReducerType:
		name = b.conf.ReducerQueueName
	}

	q, err := b.ch.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("ch.QueueDeclare: %v", err)
	}

	msgs, err := b.ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("ch.Consume: %v", err)
	}

	return msgs, nil
}

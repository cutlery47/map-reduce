package service

import (
	"fmt"

	"github.com/cutlery47/map-reduce/mapreduce"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Brocker struct {
	ch *amqp.Channel
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
		ch: ch,
	}, nil
}

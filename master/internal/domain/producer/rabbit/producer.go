package producer

import (
	"fmt"
	"io"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/mapreduce/models"
	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitTaskProducer struct {
	conn *amqp.Connection

	conf mr.Config
}

func New(conf mr.Config) (*rabbitTaskProducer, error) {
	conn, err := amqp091.Dial(fmt.Sprintf("%v:%v", conf.RabbitHost, conf.RabbitPort))
	if err != nil {
		return nil, err
	}

	return &rabbitTaskProducer{
		conn: conn,
		conf: conf,
	}, nil
}

func (rtp *rabbitTaskProducer) ProduceMapperTasks(input []io.Reader, mappers []models.Addr) ([][]byte, error) {
	return nil, nil
}

func (rtp *rabbitTaskProducer) ProduceReducerTasks(input [][]byte, reducers []models.Addr) ([][]byte, error) {
	return nil, nil
}

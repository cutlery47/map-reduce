package service

import (
	"io"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/worker/internal/domain/core"
)

type Worker struct {
	mrh  *core.MapReduceHandler
	reg  *core.RegisterHandler
	conf mr.Config
}

func New(conf mr.Config) (*Worker, error) {
	mrh, err := core.NewMapReduceHandler(conf)
	if err != nil {
		return nil, err
	}

	reg, err := core.NewRegisterHandler(conf)
	if err != nil {
		return nil, err
	}

	return &Worker{
		mrh:  mrh,
		reg:  reg,
		conf: conf,
	}, nil
}

func (w *Worker) Map(r io.Reader) (any, error) {
	return w.mrh.Map(r)
}

func (w *Worker) Reduce(mapped any) (any, error) {
	return w.mrh.Reduce(mapped)
}

func (w *Worker) Run(doneChan chan<- struct{}) error {
	// ctx, cancel := context.WithTimeout(context.TODO(), s.conf.RegisterTimeout)
	// defer cancel()

	// body := mr.WrkRegReq{
	// 	Host: s.conf.WrkHttpConf.Host,
	// 	Port: s.conf.WrkHttpConf.Port,
	// }

	// announcing worker to master and waiting for response
	// typ, err := s.r.Send(ctx, body)

	// b, err := queue.NewBrocker(conf.RabbitConf)
	// if err != nil {
	// 	errChan <- fmt.Errorf("queue.NewBrocker: %v", err)
	// 	return
	// }

	// msgs, err := b.DeclareAndConsume(regResponse.Type)
	// if err != nil {
	// 	errChan <- fmt.Errorf("b.DeclareAndConsume: %v", err)
	// 	return
	// }

	// recvChan := make(chan any)

	// go func() {
	// 	for d := range msgs {
	// 		log.Println("received a message:", string(d.Body))
	// 		reader := bytes.NewReader(d.Body)

	// 		result, err := w.Map(reader)
	// 		if err != nil {
	// 			errChan <- fmt.Errorf("w.Map: %v", err)
	// 			return
	// 		}

	// 		recvChan <- result
	// 	}

	// 	close(recvChan)
	// }()

	// for recv := range recvChan {
	// 	log.Println("handled a message:", recv)
	// }

	return nil
}

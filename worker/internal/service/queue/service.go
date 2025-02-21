package service

import (
	"io"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/worker/internal/core"
	"github.com/cutlery47/map-reduce/worker/internal/queue"
)

type Service struct {
	w *core.Worker
	r *core.Registrar

	br *queue.Brocker

	conf mr.WrkSvcConf
}

func New(conf mr.WrkSvcConf, w *core.Worker, r *core.Registrar, br *queue.Brocker) (*Service, error) {
	return &Service{
		w:    w,
		r:    r,
		br:   br,
		conf: conf,
	}, nil
}

func (s *Service) Map(r io.Reader) (any, error) {
	return s.w.Map(r)
}

func (s *Service) Reduce(mapped any) (any, error) {
	return s.w.Reduce(mapped)
}

func (s *Service) Run(doneChan chan<- struct{}) error {
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

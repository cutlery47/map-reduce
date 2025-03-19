package rabbitworker

import (
	"io"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/worker/internal/domain/core"
)

type RabbitWorker struct {
	reg  *core.RegisterHandler
	conf mr.Config
}

func New(conf mr.Config) (*RabbitWorker, error) {
	reg, err := core.NewRegisterHandler(conf)
	if err != nil {
		return nil, err
	}

	return &RabbitWorker{
		reg:  reg,
		conf: conf,
	}, nil
}

func (w *RabbitWorker) Run() error {
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

func (w *RabbitWorker) Map(input io.Reader) (io.Reader, error) {
	return mr.MapperFunc(input)
}

func (w *RabbitWorker) Reduce(input io.Reader) (io.Reader, error) {
	return mr.ReducerFunc(input)
}

func (w *RabbitWorker) Terminate(msg mr.TerminateMessage) {}

package httpworker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/worker/internal/core"
)

// core sender function decorated by http-worker required code
type HttpSenderFunction func(setupDuration, taskDuration time.Duration, errChan chan<- error, readyChan, recvChan <-chan struct{})

func UsingHttpWorker(ctx context.Context, body mapreduce.WorkerRegisterRequest, send core.SenderFunction) HttpSenderFunction {
	return func(setupDuration, taskDuration time.Duration, errChan chan<- error, readyChan, recvChan <-chan struct{}) {
		// waiting for http server to set up
		readyTimer := time.NewTimer(setupDuration)
		select {
		case <-readyTimer.C:
			errChan <- fmt.Errorf("server setup timed out")
			return
		case <-readyChan:
		}

		send(ctx, body)

		// waiting for new jobs
		// terminate worker if no jobs were received
		go func() {
			timer := time.NewTimer(taskDuration)
			log.Println("waiting for tasks for:", taskDuration)

			select {
			case <-recvChan:
				return
			case <-timer.C:
				errChan <- fmt.Errorf("didn't receive any tasks")
				return
			}
		}()
	}
}

package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	sv              *http.Server
	listener        net.Listener
	shutdownTimeout time.Duration
}

func New(conf mr.Config, listener net.Listener, handler http.Handler) *Server {
	httpserv := &http.Server{
		Handler:      handler,
		ReadTimeout:  conf.WorkerReadTimeout,
		WriteTimeout: conf.WorkerWriteTimeout,
		Addr:         fmt.Sprintf("%v:%v", conf.WorkerHost, listener.Addr().(*net.TCPAddr).Port),
	}

	return &Server{
		sv:              httpserv,
		listener:        listener,
		shutdownTimeout: conf.WorkerShutdownTimeout,
	}
}

func (s *Server) Run(doneCh <-chan error) error {
	log.Infoln("[HTTP-SERVER] running http server on:", s.sv.Addr)

	go func() {
		if err := s.sv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Println("[HTTP-SERVER] error:", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// waiting for either kernel signal or app signal
	select {
	case <-sigChan:
		// received kernel signal
	case err := <-doneCh:
		// master has finished execution
		if err != nil {
			log.Errorln("[RUNTIME ERROR] error:", err)
		} else {
			log.Infoln("[MASTER] finishing execution...")
		}
	}

	log.Infoln("[HTTP-SERVER] shutting down gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.sv.Shutdown(ctx)
}

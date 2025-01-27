package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultAdress          = "0.0.0.0:0"
	defaultReadTimeout     = 3 * time.Second
	defaultWriteTimeout    = 3 * time.Second
	defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	Server *http.Server

	listener net.Listener

	readyChan chan<- struct{}
	errChan   <-chan error
	endChan   <-chan struct{}

	shutdownTimeout time.Duration
}

func New(handler http.Handler, listener net.Listener, readyChan chan<- struct{}, errChan <-chan error, endChan <-chan struct{}) *Server {
	httpserv := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		Addr:         defaultAdress,
	}

	serv := &Server{
		Server:          httpserv,
		listener:        listener,
		readyChan:       readyChan,
		errChan:         errChan,
		endChan:         endChan,
		shutdownTimeout: defaultShutdownTimeout,
	}

	return serv
}

func (s *Server) Run(ctx context.Context) error {
	log.Println(fmt.Sprintf("running http Server on %v", s.Server.Addr))

	go func() {
		if err := s.Server.Serve(s.listener); !errors.Is(err, http.ErrServerClosed) {
			log.Println("http Server error:", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// signaling that server has been set up
	s.readyChan <- struct{}{}

	// waiting for either kernel signal or app signal
	select {
	case <-s.endChan:
	case <-sigChan:
	case err := <-s.errChan:
		log.Println("error:", err)
	}

	log.Println("shutting down http-server")

	ctx, cancel := context.WithTimeout(ctx, s.shutdownTimeout)
	defer cancel()

	return s.Server.Shutdown(ctx)
}

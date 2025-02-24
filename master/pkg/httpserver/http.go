package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	mr "github.com/cutlery47/map-reduce/mapreduce"
)

const (
	defaultAdress          = "0.0.0.0:8080"
	defaultReadTimeout     = 3 * time.Second
	defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	sv *http.Server

	shutdownTimeout time.Duration
}

func New(conf mr.MstHttpConf, handler http.Handler) *Server {
	httpserv := &http.Server{
		Handler:     handler,
		ReadTimeout: defaultReadTimeout,
		Addr:        fmt.Sprintf("%v:%v", conf.Host, conf.Port),
	}

	serv := &Server{
		sv:              httpserv,
		shutdownTimeout: defaultShutdownTimeout,
	}

	return serv
}

func (s *Server) Run(doneChan <-chan struct{}, errChan <-chan error) error {
	log.Println(fmt.Sprintf("running http server on %v", s.sv.Addr))

	go func() {
		if err := s.sv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Println("http server error:", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// waiting for either kernel signal or app signal
	select {
	case <-doneChan:
	case <-sigChan:
	case err := <-errChan:
		log.Println("error:", err)
	}

	log.Println("shutting down http-server")

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.sv.Shutdown(ctx)
}

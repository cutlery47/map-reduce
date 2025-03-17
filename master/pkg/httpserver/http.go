package httpserver

import (
	"context"
	"errors"
	"fmt"
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
	shutdownTimeout time.Duration
}

func New(conf mr.Config, handler http.Handler) *Server {
	httpserv := &http.Server{
		Handler:      handler,
		ReadTimeout:  conf.MasterReadTimeout,
		WriteTimeout: conf.MasterWriteTimeout,
		Addr:         fmt.Sprintf("%v:%v", conf.MasterHost, conf.MasterPort),
	}

	serv := &Server{
		sv:              httpserv,
		shutdownTimeout: conf.MasterShutdownTimeout,
	}

	return serv
}

func (s *Server) Run(doneChan <-chan struct{}, errChan <-chan error) error {
	log.Infoln("[HTTP-SERVER] running http server on:", s.sv.Addr)

	go func() {
		if err := s.sv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Errorln("[HTTP-SERVER] error:", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		// received kernel signal
	case <-doneChan:
		// finished gracefully
	case err := <-errChan:
		// received error
		log.Errorln("[RUNTIME ERROR] error:", err)
	}

	log.Infoln("[HTTP-SERVER] shutting down gracefully")
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.sv.Shutdown(ctx)
}

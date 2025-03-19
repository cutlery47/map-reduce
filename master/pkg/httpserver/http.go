package httpserver

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	defaultAddr            = "127.0.0.1:8080"
	defaultReadTimeout     = 3 * time.Second
	defaultWriteTimeout    = 3 * time.Second
	defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	hs              *http.Server
	shutdownTimeout time.Duration
}

func New(handler http.Handler, opts ...Option) *Server {
	hs := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		Addr:         defaultAddr,
	}

	s := &Server{
		hs:              hs,
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Run(doneCh <-chan AppSignal) error {
	log.Infoln("[HTTP-SERVER] running http server on:", s.hs.Addr)

	go func() {
		if err := s.hs.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Errorln("[HTTP-SERVER] error:", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// waiting for either kernel signal or app signal
	select {
	case <-sigChan:
		// received kernel signal
	case sig := <-doneCh:
		// received application signal
		if sig.Error != nil {
			log.Errorln("[RUNTIME ERROR] error:", sig.Error)
		} else {
			log.Infoln("[MASTER] finished with message:", sig.Message)
		}
	}

	log.Infoln("[HTTP-SERVER] shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.hs.Shutdown(ctx)
}

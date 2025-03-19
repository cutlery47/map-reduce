package httpserver

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	defaultReadTimeout     = 3 * time.Second
	defaultWriteTimeout    = 3 * time.Second
	defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	hs              *http.Server
	listener        net.Listener
	shutdownTimeout time.Duration
}

func New(handler http.Handler, listener net.Listener, opts ...Option) *Server {
	hs := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	s := &Server{
		hs:              hs,
		listener:        listener,
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Run(doneCh <-chan AppSignal) error {
	log.Infoln("[HTTP-SERVER] running http server on:", s.listener.Addr())

	go func() {
		if err := s.hs.Serve(s.listener); !errors.Is(err, http.ErrServerClosed) {
			log.Println("[HTTP-SERVER] error:", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// waiting for either kernel signal or app signal
	select {
	case <-sigChan:
		// received kernel signal
	case sig := <-doneCh:
		// worker has finished execution
		if sig.Error != nil {
			log.Errorln("[RUNTIME ERROR] error:", sig.Error)
		} else {
			log.Infoln("[HTTP-SERVER] finishing execution...")
		}
	}

	log.Infoln("[HTTP-SERVER] shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.hs.Shutdown(ctx)
}

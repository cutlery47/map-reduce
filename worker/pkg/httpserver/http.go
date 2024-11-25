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
)

const (
	defaultAdress          = "0.0.0.0:8080"
	defaultReadTimeout     = 3 * time.Second
	defaultWriteTimeout    = 3 * time.Second
	defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	server *http.Server

	shutdownTimeout time.Duration
}

func New(handler http.Handler) *Server {
	httpserv := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		Addr:         defaultAdress,
	}

	serv := &Server{
		server:          httpserv,
		shutdownTimeout: defaultShutdownTimeout,
	}

	return serv
}

func (s *Server) Run(ctx context.Context) error {
	log.Println(fmt.Sprintf("running http server on %v", s.server.Addr))

	go func() {
		if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Println("http server error:", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	log.Println("Shutting down http server")

	ctx, cancel := context.WithTimeout(ctx, s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}

package httpserver

import "time"

type Option func(s *Server)

func WithAddr(addr string) Option {
	return func(s *Server) {
		s.hs.Addr = addr
	}
}

func WithReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.hs.ReadTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.hs.WriteTimeout = timeout
	}
}

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}

package httpserver

import (
	"io"
	"net/http"
)

type Server struct {
	srv *http.Server
	cli *http.Client

	// channel for receiving requests
	reqChan chan<- Request
	// channel for passing responses
	resChan <-chan Response
	// channel for passing errors
	errChan chan<- error
}

type Request struct {
	Addr string
	Body io.Reader
}

type Response struct {
	Addr string
	Body io.Reader
}

func New(reqChan chan<- Request, resChan <-chan Response, errChan chan<- error) *Server {
	srv := &http.Server{
		Addr: "0.0.0.0:8080",
	}

	cli := &http.Client{}

	return &Server{
		srv:     srv,
		cli:     cli,
		reqChan: reqChan,
		resChan: resChan,
		errChan: errChan,
	}
}

func (s *Server) ListenForRequests() {

}

func (s *Server) ListenFromResponses() {

}

// intended to be run in a goroutine
func (s *Server) Run() {
	err := s.srv.ListenAndServe()
	if err != nil {
		s.errChan <- err
	}
}

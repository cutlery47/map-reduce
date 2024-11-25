package main

import (
	"context"
	"net/http"

	"github.com/cutlery47/map-reduce/master/pkg/httpserver"
)

type Handler struct {
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello, worker!"))
}

func main() {
	ctx := context.Background()

	mux := http.NewServeMux()
	mux.Handle("/", Handler{})

	server := httpserver.New(mux)
	server.Run(ctx)
}

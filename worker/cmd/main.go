package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/worker/pkg/httpserver"
)

var (
	mapFunc, reduceFunc = mapreduce.MapReduce()
)

func SendRegister(ctx context.Context, cl *http.Client, port int, readyChan <-chan struct{}, errChan chan<- error, readyTimeout, requestTimeout time.Duration) {
	readyTimer := time.NewTimer(readyTimeout)

	select {
	case <-readyTimer.C:
		errChan <- fmt.Errorf("server setup timed out")
		return
	case <-readyChan:
	}

	body := mapreduce.WorkerRegisterRequest{
		Port: strconv.Itoa(port),
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		errChan <- fmt.Errorf("error: %v", err)
		return
	}
	jsonReader := bytes.NewReader(jsonBody)

	reqCtx, cancel := context.WithTimeout(ctx, readyTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, "POST", "http://localhost:8080/register", jsonReader)
	if err != nil {
		errChan <- fmt.Errorf("http.NewRequestWithContext: %v", err)
		return
	}

	res, err := cl.Do(req)
	if err != nil {
		errChan <- fmt.Errorf("http.Get: %v", err)
		return
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		errChan <- fmt.Errorf("io.ReadAll: %v", err)
		return
	}

	log.Println("response:", string(resBody))
}

func main() {
	ctx := context.Background()
	cl := http.DefaultClient
	cl.Timeout = 3 * time.Second
	log.Println("started client")

	readyChan := make(chan struct{})
	errChan := make(chan error)

	http.HandleFunc("/map", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(400)
			w.Write([]byte("not found"))
			return
		}

		res, err := mapFunc(r.Body)
		if err != nil {
			errChan <- fmt.Errorf("mapFunc: %v", err)
			w.WriteHeader(500)
			w.Write([]byte("internal server error"))
			return
		}

		json, err := json.Marshal(res)
		if err != nil {
			errChan <- fmt.Errorf("json.Marshall: %v", err)
			w.WriteHeader(500)
			w.Write([]byte("internal server error"))
			return
		}

		w.Write(json)
		w.WriteHeader(200)
	})

	http.HandleFunc("/reduce", func(w http.ResponseWriter, r *http.Request) {

	})

	// creating a new http server
	srv := httpserver.New(http.DefaultServeMux)

	// meanwhile running the http server
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("net.Listen: ", err)
	}

	go SendRegister(ctx, cl, listener.Addr().(*net.TCPAddr).Port, readyChan, errChan, time.Second*3, time.Second*3)
	log.Fatal("server error: ", srv.Run(ctx, listener, readyChan, errChan))
}

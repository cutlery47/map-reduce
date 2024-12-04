package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/master/pkg/httpserver"
)

type Addr struct {
	Port string
	Host string
}

var (
	workerAddrs  []Addr
	mapperAddrs  []Addr
	reducerAddrs []Addr
)

func DistWorkers(cl *http.Client, mappers, reducers int) {
	for len(workerAddrs) != mappers+reducers {
		log.Println(workerAddrs)
		time.Sleep(2 * time.Second)
	}

	for _, addr := range workerAddrs {
		req := mapreduce.MasterDistrRequest{}
		res := &mapreduce.WorkerAckResponse{}

		if mappers > 0 {
			req.Distr = mapreduce.MapperDistr
			mappers--
		} else if reducers > 0 {
			req.Distr = mapreduce.ReducerDistr
			reducers--
		} else {
			panic("received more worker requests than intended")
		}

		data, err := json.Marshal(req)
		if err != nil {
			log.Println(err)
		}

		body := bytes.NewReader(data)
		httpRes, err := cl.Post(fmt.Sprintf("http://%v:%v/dist", addr.Host, addr.Port), "application/json", body)
		if err != nil {
			log.Println(err)
		}

		resBody, err := io.ReadAll(httpRes.Body)
		if err != nil {
			log.Println(err)
		}

		err = json.Unmarshal(resBody, res)
		if err != nil {
			log.Println(err)
		}

		if res.Ack != mapreduce.WorkerAck {
			log.Println("worker ack not returned")
		}
	}
}

func main() {
	ctx := context.Background()
	mu := &sync.Mutex{}

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(400)
			w.Write([]byte("not found"))
		}

		body := r.Body
		req := mapreduce.WorkerRegisterRequest{}

		jsonBody, err := io.ReadAll(body)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("internal server error: " + err.Error()))
			return
		}

		err = json.Unmarshal(jsonBody, &req)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("internal server error: " + err.Error()))
			return
		}

		addr := Addr{
			Port: req.Port,
			Host: strings.Split(r.RemoteAddr, ":")[0],
		}

		mu.Lock()
		workerAddrs = append(workerAddrs, addr)
		mu.Unlock()

		w.Write([]byte("hello, worker!"))
	})

	go DistWorkers(http.DefaultClient, 1, 1)

	server := httpserver.New(http.DefaultServeMux)
	log.Fatal("error: ", server.Run(ctx))
}

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/master/pkg/httpserver"
)

type Addr struct {
	Port string
	Host string
}

var (
	counter      atomic.Int64
	mapperAddrs  []Addr
	reducerAddrs []Addr
)

var (
	mapMu = &sync.Mutex{}
	redMu = &sync.Mutex{}
)

func CollectWorkers(cl *http.Client, mappers, reducers int, timeout time.Duration, errChan chan<- error) {
	timer := time.NewTimer(timeout)

	for cnt := counter.Load(); cnt != int64(mappers)+int64(reducers); cnt = counter.Load() {
		select {
		case <-timer.C:
			errChan <- fmt.Errorf("expected %v workers, connected: %v", mappers+reducers, cnt)
			return
		default:
			log.Println("connected:", cnt)
			time.Sleep(time.Second)
		}
	}

	log.Println("all workers connected")
	go AssignMappers(cl, errChan)
}

func AssignMappers(cl *http.Client, errChan chan<- error) {
	files, err := SplitFile("file.txt", len(mapperAddrs))
	if err != nil {
		errChan <- fmt.Errorf("SplitFile: %v", err)
		return
	}

	var mapResults [][]byte
	var redResults [][]byte

	for i, addr := range mapperAddrs {
		res, err := cl.Post(fmt.Sprintf("http://%v:%v/map", addr.Host, addr.Port), "text/plain", files[i])
		if err != nil {
			errChan <- err
			return
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			errChan <- err
			return
		}

		mapResults = append(mapResults, body)
	}

	for i, addr := range reducerAddrs {
		res, err := cl.Post(fmt.Sprintf("http://%v:%v/reduce", addr.Host, addr.Port), "application/json", bytes.NewReader(mapResults[i]))
		if err != nil {
			errChan <- err
			return
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			errChan <- err
			return
		}

		redResults = append(redResults, body)
	}

	if err := ReturnResults(redResults); err != nil {
		errChan <- err
	}
}

func ReturnResults(redResults [][]byte) error {
	hashMap := make(map[string]int)

	for _, res := range redResults {
		sliceRes := []string{}
		json.Unmarshal(res, &sliceRes)
		fmt.Println(sliceRes)

		for _, el := range sliceRes {
			elSplit := strings.Split(el, ":")
			key := elSplit[0]
			val, err := strconv.Atoi(elSplit[1])
			if err != nil {
				return err
			}

			if v, ok := hashMap[key]; !ok {
				hashMap[key] = val
			} else {
				hashMap[key] = v + val
			}
		}
	}

	fmt.Println(hashMap)

	return nil
}

func SplitFile(filename string, mappers int) ([]io.Reader, error) {
	readers := []io.Reader{}

	err := os.Mkdir("chunks", 0777)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			return nil, fmt.Errorf("os.Mkdir: %v", err)
		}
	}

	bash := "split"
	arg0, arg1 := "-n", strconv.Itoa(mappers)
	arg2 := "-d"
	arg3, arg4 := "file.txt", "chunks/"

	cmd := exec.Command(bash, arg0, arg1, arg2, arg3, arg4)
	cmd.Stderr = os.Stdin
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	chunks, err := os.ReadDir("chunks")
	if err != nil {
		return nil, fmt.Errorf("os.ReadDir: %v", err)
	}

	for _, chunk := range chunks {
		fd, err := os.Open(fmt.Sprintf("chunks/%v", chunk.Name()))
		if err != nil {
			return nil, err
		}

		readers = append(readers, fd)
	}

	return readers, nil
}

func main() {
	ctx := context.Background()
	syscall.Umask(0)

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(400)
			w.Write([]byte("not found"))
			return
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

		cur := counter.Add(1)
		if cur%2 == 0 {
			mapMu.Lock()
			mapperAddrs = append(mapperAddrs, addr)
			mapMu.Unlock()
		} else {
			redMu.Lock()
			reducerAddrs = append(reducerAddrs, addr)
			redMu.Unlock()
		}

		w.Write([]byte("hello, worker!"))
	})

	errChan := make(chan error)

	go CollectWorkers(http.DefaultClient, 1, 1, time.Second*10, errChan)

	server := httpserver.New(http.DefaultServeMux)
	log.Fatal("server error: ", server.Run(ctx, errChan))
}

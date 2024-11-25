package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	res, err := http.Get("http://master-service:8080/")
	if err != nil {
		log.Fatal("http.Get:", err)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("io.ReadAll:", err)
	}

	log.Println("response:", string(resBody))
}

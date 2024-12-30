package main

import (
	"log"

	"github.com/cutlery47/map-reduce/worker/internal/app"
)

func main() {
	log.Fatal("error: ", app.Run())
}

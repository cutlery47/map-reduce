package main

import (
	"flag"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/worker/internal/app"
	log "github.com/sirupsen/logrus"
)

var envLocation = flag.String("env", ".env", "specify env-file name and location")

func main() {
	// parse input parameters
	flag.Parse()

	conf, err := mr.NewConfig(*envLocation)
	if err != nil {
		log.Fatalf("[SETUP] error reading config: %v", err)
	}

	err = app.Run(*conf)
	if err != nil {
		log.Fatalf("[SETUP] error setting up worker: %v", err)
	}
}

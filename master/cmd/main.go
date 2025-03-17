package main

import (
	"flag"
	"syscall"

	mr "github.com/cutlery47/map-reduce/mapreduce"
	"github.com/cutlery47/map-reduce/master/internal/app"
	log "github.com/sirupsen/logrus"
)

var envLocation = flag.String("env", ".env", "specify env-file name and location")

func main() {
	// drop permission limit
	syscall.Umask(0)
	// parse input parameters
	flag.Parse()

	conf, err := mr.NewConfig(*envLocation)
	if err != nil {
		log.Fatalf("[SETUP] error when reading config: %v", err)
	}

	app.Run(*conf)
}

package main

import (
	"log"
	"time"

	"github.com/cutlery47/mapreduce"
)

func main() {
	mapFunc, reduceFunc := mapreduce.MapReduce()

	log.Println(mapFunc)
	log.Println(reduceFunc)

	time.Sleep(time.Hour)
}

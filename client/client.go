package client

import (
	"log"
	"os"

	mapreduce "github.com/cutlery47/map-reduce"
)

func Run() {
	mapFunc, reduceFunc := mapreduce.MapReduce()

	file, err := os.Open("file.txt")
	if err != nil {
		log.Fatal(err)
	}

	mapResult, err := mapFunc(file)
	if err != nil {
		log.Fatal("Map:", err)
	}

	reduceResult, err := reduceFunc(mapResult)
	if err != nil {
		log.Fatal("Reduce:", err)
	}

	log.Println(reduceResult)
}

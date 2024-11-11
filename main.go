package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {

	var mapFunc Map[string, int] = func(input io.Reader) ([]MapResult[string, int], error) {
		res := []MapResult[string, int]{}

		raw, err := io.ReadAll(input)
		if err != nil {
			return nil, err
		}

		text := string(raw)

		words := strings.Split(strip(text), " ")
		for _, word := range words {
			res = append(res, MapResult[string, int]{
				MappedKey:   word,
				MappedValue: 1,
			})
		}

		return res, nil
	}

	var reduceFunc Reduce[string, int, string] = func(res []MapResult[string, int]) ([]string, error) {
		reducerRes := []string{}

		freq := map[string]int{}

		for _, mres := range res {
			if v, ok := freq[mres.MappedKey]; !ok {
				freq[mres.MappedKey] = 1
			} else {
				freq[mres.MappedKey] = v + 1
			}
		}

		for k, v := range freq {
			reducerRes = append(reducerRes, fmt.Sprintf("%v:%v", k, v))
		}

		return reducerRes, nil
	}

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

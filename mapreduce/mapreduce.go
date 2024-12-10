package mapreduce

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/cutlery47/map-reduce/mapreduce/internal"
)

// KT - Map result key type
// VT - Map result value type
// RT - Reduce result type

type Map[KT, VT any] func(input io.Reader) ([]MapResult[KT, VT], error)

type MapResult[KT, VT any] struct {
	MappedKey   KT `json:"key"`
	MappedValue VT `json:"value"`
}

type Reduce[KT, VT, RT any] func(res []MapResult[KT, VT]) ([]RT, error)

// Function which returns user implementations of Map() and Reduce()
func MapReduce() (Map[string, int], Reduce[string, int, string]) {
	return MyMap, MyReduce
}

// Define your Map(), Reduce(), MapResult and RedResult implementations below...
// ======================================================

type MyMapResult []MapResult[string, int]
type MyRedResult []string

var MyMap Map[string, int] = func(input io.Reader) ([]MapResult[string, int], error) {
	log.Println("in map")

	res := []MapResult[string, int]{}

	raw, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	text := string(raw)

	words := strings.Split(internal.Strip(text), " ")
	for _, word := range words {
		res = append(res, MapResult[string, int]{
			MappedKey:   word,
			MappedValue: 1,
		})
	}

	return res, nil
}

var MyReduce Reduce[string, int, string] = func(res []MapResult[string, int]) ([]string, error) {
	log.Println("in reduce")

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

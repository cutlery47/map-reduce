package mapreduce

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

// Interface for your personal mapper / reducer functions
type Func func(input io.Reader) (io.Reader, error)

// Define your Func implementations below
// ===============================================
// ATTENTION: IT IS NECESSARY TO NAME YOUR FUNCTIONS "MapperFunc" AND "ReducerFunc" RESPECTIVELY

// Code below is for example purposes

type Entry struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

type (
	MapResult []Entry
	RedResult []Entry
)

var MapperFunc Func = func(input io.Reader) (io.Reader, error) {
	var (
		res MapResult
	)

	raw, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	for _, word := range strings.Split(string(raw), " ") {
		res = append(res, Entry{
			Key:   word,
			Value: 1,
		})
	}

	resJson, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(resJson), nil
}

var ReducerFunc Func = func(input io.Reader) (io.Reader, error) {
	var (
		res RedResult
	)

	err := json.NewDecoder(input).Decode(&res)
	if err != nil {
		return nil, err
	}

	resJson, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(resJson), nil
}

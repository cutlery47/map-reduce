package mapreduce

import "time"

type Config struct {
	Mappers  int
	Reducers int

	RegisterTimeout time.Duration
	CollectTimeout  time.Duration
	ReadyTimeout    time.Duration
	RequestTimeout  time.Duration
}

var DefaultConfig = &Config{
	Mappers:  1,
	Reducers: 1,

	RegisterTimeout: 10 * time.Second,
	CollectTimeout:  1 * time.Second,
	ReadyTimeout:    5 * time.Second,
	RequestTimeout:  3 * time.Second,
}

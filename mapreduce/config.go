package mapreduce

type Config struct {
	Mappers  int
	Reducers int
}

var DefaultConfig = &Config{
	Mappers:  1,
	Reducers: 1,
}

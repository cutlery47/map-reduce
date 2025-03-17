package producer

import (
	"io"

	mr "github.com/cutlery47/map-reduce/mapreduce"
)

type TaskProducer interface {
	ProduceMapperTasks(input []io.Reader, mappers []mr.Addr) ([][]byte, error)
	ProduceReducerTasks(input [][]byte, reducers []mr.Addr) ([][]byte, error)
}

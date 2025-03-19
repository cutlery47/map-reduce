package producer

import (
	"io"

	"github.com/cutlery47/map-reduce/mapreduce/models"
)

type TaskProducer interface {
	ProduceMapperTasks(input []io.Reader, mappers []models.Addr) ([][]byte, error)
	ProduceReducerTasks(input [][]byte, reducers []models.Addr) ([][]byte, error)
}

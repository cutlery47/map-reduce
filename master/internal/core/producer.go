package core

import (
	"io"
)

type Addr struct {
	Port string
	Host string
}

type TaskProducer interface {
	ProduceMapperTasks(files []io.Reader) ([][]byte, error)
	ProduceReducerTasks(mapResults [][]byte) ([][]byte, error)
}

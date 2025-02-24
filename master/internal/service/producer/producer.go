package core

import (
	"io"
)

type TaskProducer interface {
	ProduceMapperTasks(in []io.Reader) ([][]byte, error)
	ProduceReducerTasks(in [][]byte) ([][]byte, error)
}

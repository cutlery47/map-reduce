package worker

import (
	"io"

	"github.com/cutlery47/map-reduce/mapreduce/models"
)

// general interface for any worker
type Worker interface {
	// main worker loop
	Run() error
	// starts map procedure
	Map(input io.Reader) (io.Reader, error)
	// starts reduce procedure
	Reduce(input io.Reader) (io.Reader, error)
	// terminates worker
	Terminate(msg models.TerminateMessage)
}

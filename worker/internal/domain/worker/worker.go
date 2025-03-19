package worker

import (
	"io"

	mr "github.com/cutlery47/map-reduce/mapreduce"
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
	Terminate(msg mr.TerminateMessage)
}

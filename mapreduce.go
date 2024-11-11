package main

import "io"

// A worker who is assigned a map task reads the
// contents of the corresponding input split. It parses
// key/value pairs out of the input data and passes each
// pair to the user-defined Map function.
// The intermediate key/value pairs produced by the Map function
// are buffered in memory.

// KT - Result key type
// VT - Result value type

type Map[KT, VT any] func(input io.Reader) ([]MapResult[KT, VT], error)

type MapResult[KT, VT any] struct {
	MappedKey   KT
	MappedValue VT
}

// Periodically, the buffered pairs are written to local
// disk, partitioned into R regions by the partitioning
// function. The locations of these buffered pairs on
// the local disk are passed back to the master, who
// is responsible for forwarding these locations to the
// reduce workers.

// When a reduce worker is notified by the master
// about these locations, it uses remote procedure calls
// to read the buffered data from the local disks of the
// map workers. When a reduce worker has read all intermediate data, it sorts it by the intermediate keys
// so that all occurrences of the same key are grouped
// together. The sorting is needed because typically
// many different keys map to the same reduce task. If
// the amount of intermediate data is too large to fit in
// memory, an external sort is used.

// KT - MapResult key type
// VT - MapResult value type
// RT - Reduce result type

type Reduce[KT, VT, RT any] func(res []MapResult[KT, VT]) ([]RT, error)

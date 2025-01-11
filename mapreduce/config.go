package mapreduce

import "time"

type Config struct {
	// master host
	MasterHost string
	// master pod
	MasterPort string

	// amount of mappers
	Mappers int
	// amount of reducers
	Reducers int

	// location of a file to be mapped
	FileLocation string
	// location of chunksn
	ChunkLocation string
	// location of results
	ResultLocation string

	// time for all worker nodes to register on master
	RegisterDuration time.Duration
	// time in between which master node checks for registered worker nodes
	CollectTimeout time.Duration
	// time for worker node to set up
	ReadyTimeout time.Duration
	// time for worker to receive a new job
	WorkerTimeout time.Duration
}

var LocalConfig = Config{
	MasterHost: "localhost",
	MasterPort: "8080",

	Mappers:  1,
	Reducers: 1,

	FileLocation:   "$HOME/mapreduce/file.txt",
	ChunkLocation:  "$HOME/mapreduce/chunks",
	ResultLocation: "$HOME/mapreduce/results",

	RegisterDuration: 10 * time.Second,
	CollectTimeout:   1 * time.Second,
	ReadyTimeout:     5 * time.Second,
	WorkerTimeout:    10 * time.Second,
}

var KubernetesConfig = Config{
	MasterHost: "master-service",
	MasterPort: "8080",

	Mappers:  1,
	Reducers: 1,

	FileLocation:   "/files/file.txt",
	ChunkLocation:  "/files/chunks",
	ResultLocation: "/files/results",

	RegisterDuration: 10 * time.Second,
	CollectTimeout:   1 * time.Second,
	ReadyTimeout:     5 * time.Second,
	WorkerTimeout:    10 * time.Second,
}

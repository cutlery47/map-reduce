package mapreduce

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	// master host
	MasterHost string `env:"MASTER_HOST"`
	// master pod
	MasterPort string `env:"MASTER_PORT"`

	// amount of mappers
	Mappers int `env:"MAPPERS"`
	// amount of reducers
	Reducers int `env:"REDUCERS"`

	// name of a file to be mapped
	FileName string `env:"FILE_NAME"`
	// directory to be mapped into master node
	FileDirectory string `env:"FILE_DIRECTORY"`
	// location of chunks inside file directory
	ChunkDirectory string `env:"CHUNK_DIRECTORY"`
	// location of results inside file directory
	ResultDirectory string `env:"RESULT_DIRECTORY"`

	// time for all worker nodes to register on master
	RegisterDuration time.Duration `env:"REGISTER_DURATION"`
	// time for worker node to set up
	SetupDuration time.Duration `env:"SETUP_DURATION"`
	// time in between which master node checks for registered worker nodes
	CollectInterval time.Duration `env:"COLLECT_INTERVAL"`
	// time for worker to receive a new job
	WorkerTimeout time.Duration `env:"WORKER_TIMEOUT"`
}

func NewConfig(envLocation string) (Config, error) {
	conf := Config{}

	if err := godotenv.Load(envLocation); err != nil {
		return conf, err
	}

	if err := cleanenv.ReadEnv(&conf); err != nil {
		return conf, err
	}

	return conf, nil
}

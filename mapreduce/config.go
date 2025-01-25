package mapreduce

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type WorkerConfig struct {
	WorkerRegistrarConfig

	ProducerType string `env:"PRODUCER_TYPE"`
}

type WorkerRegistrarConfig struct {
	MasterHost string `env:"MASTER_HOST"`
	MasterPort string `env:"MASTER_PORT"`

	// time for worker node to set up
	SetupDuration time.Duration `env:"SETUP_DURATION"`
	// time for worker to receive a new job
	WorkerTimeout time.Duration `env:"WORKER_TIMEOUT"`
}

type MasterConfig struct {
	MasterRegistrarConfig
	ProducerConfig

	ProducerType string `env:"PRODUCER_TYPE"`

	// directory to be mapped into master node
	FileDirectory string `env:"MAPPED_FILE_DIRECTORY"`
	// name of a file to be mapped
	FileName string `env:"FILE_NAME"`
	// location of chunks inside file directory
	ChunkDirectory string `env:"CHUNK_DIRECTORY"`
	// location of results inside file directory
	ResultDirectory string `env:"RESULT_DIRECTORY"`
}

type MasterRegistrarConfig struct {
	// amount of mappers
	Mappers int `env:"MAPPERS"`
	// amount of reducers
	Reducers int `env:"REDUCERS"`

	// time for all worker nodes to register on master
	RegisterDuration time.Duration `env:"REGISTER_DURATION"`
	// time in between which master node checks for registered worker nodes
	CollectInterval time.Duration `env:"COLLECT_INTERVAL"`
}

type ProducerConfig struct {
	RabbitMapperPath  string `env:"RABBIT_MAPPER_PATH"`
	RabbitReducerPath string `env:"RABBIT_REDUCER_PATH"`
	RabbitConfig
}

type RabbitConfig struct {
	RabbitLogin    string `env:"RABBIT_LOGIN"`
	RabbitPassword string `env:"RABBIT_PASSWORD"`
	RabbitHost     string `env:"RABBIT_HOST"`
	RabbitPort     string `env:"RABBIT_PORT"`
}

func NewMasterConfig(envLocation string) (MasterConfig, error) {
	conf := MasterConfig{}

	if err := readConfig(envLocation, &conf); err != nil {
		return conf, fmt.Errorf("readConfig: %v", err)
	}

	return conf, nil
}

func NewWorkerConfig(envLocation string) (WorkerConfig, error) {
	conf := WorkerConfig{}

	if err := readConfig(envLocation, &conf); err != nil {
		return conf, fmt.Errorf("readConfig: %v", err)
	}

	return conf, nil
}

func NewRabbitConfig(envLocation string) (RabbitConfig, error) {
	conf := RabbitConfig{}

	if err := readConfig(envLocation, &conf); err != nil {
		return conf, fmt.Errorf("readConfig: %v", err)
	}

	return conf, nil
}

func readConfig(envLocation string, conf any) error {
	if err := godotenv.Overload(envLocation); err != nil {
		return fmt.Errorf("godotenv.Load: %v", err)
	}

	if err := cleanenv.ReadEnv(conf); err != nil {
		return fmt.Errorf("cleanenv.ReadEnv: %v", err)
	}

	return nil
}

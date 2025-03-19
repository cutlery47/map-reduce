package mapreduce

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	// ================== GENERAL ==================
	// file to be processed
	File string `env:"FILE" env-required:""`
	// amount of workers
	Mappers int `env:"MAPPERS" env-required:""`
	// amount of reducers
	Reducers int `env:"REDUCERS" env-required:""`
	// transport to be used in the system
	Transport string `env:"TRANSPORT" env-required:""`
	// directory to map into master node
	MappedDir string `env:"MAPPED_DIRECTORY" env-default:"data"`
	// time for all worker nodes to register on master
	RegisterDur time.Duration `env:"REGISTER_DURATION" env-default:"30s"`

	// ================== MASTER ==================
	MasterHost            string        `env:"MASTER_HOST"              env-default:"localhost"`
	MasterPort            string        `env:"MASTER_PORT"              env-default:"8080"`
	MasterReadTimeout     time.Duration `env:"MASTER_READ_TIMEOUT"      env-default:"3s"`
	MasterWriteTimeout    time.Duration `env:"MASTER_WRITE_TIMEOUT"     env-default:"3s"`
	MasterShutdownTimeout time.Duration `env:"MASTER_SHUTDOWN_TIMEOUT"  env-default:"3s"`
	// time for master to receive request from any worker
	MasterRequestAwaitDur time.Duration `env:"MASTER_REQUEST_AWAIT_DURATION" env-default:"10s"`

	// ================== WORKER ==================
	WorkerHost            string        `env:"WORKER_HOST"              env-default:"localhost"`
	WorkerReadTimeout     time.Duration `env:"WORKER_READ_TIMEOUT"      env-default:"3s"`
	WorkerWriteTimeout    time.Duration `env:"WORKER_WRITE_TIMEOUT"     env-default:"3s"`
	WorkerShutdownTimeout time.Duration `env:"WORKER_SHUTDOWN_TIMEOUT"  env-default:"3s"`
	// time for worker node to set up
	WorkerSetupDur time.Duration `env:"WORKER_SETUP_DURATION" env-default:"10s"`
	// time for worker to receive a new job
	WorkerAwaitDur time.Duration `env:"WORKER_AWAIT_DURATION" env-default:"10s"`

	// ================== RABBITMQ ==================
	RabbitLogin        string `env:"RABBIT_LOGIN"         env-default:"login"`
	RabbitPassword     string `env:"RABBIT_PASSWORD"      env-defualt:"password"`
	RabbitHost         string `env:"RABBIT_HOST"          env-default:"localhost"`
	RabbitPort         string `env:"RABBIT_PORT"          env-default:"5673"`
	RabbitMapperQueue  string `env:"RABBIT_MAPPER_QUEUE"  env-default:"mapper"`
	RabbitReducerQueue string `env:"RABBIT_REDUCER_QUEUE" env-default:"reducer"`
}

func NewConfig(envLocation string) (*Config, error) {
	var conf Config

	if err := godotenv.Overload(envLocation); err != nil {
		return nil, fmt.Errorf("godotenv.Load: %v", err)
	}

	if err := cleanenv.ReadEnv(&conf); err != nil {
		return nil, fmt.Errorf("cleanenv.ReadEnv: %v", err)
	}

	return &conf, nil
}

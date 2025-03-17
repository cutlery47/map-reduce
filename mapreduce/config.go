package mapreduce

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// ================== WORKER CONFIG ==================

type Config struct {
	// ================== GENERAL ==================
	// file to be processed
	File string `env:"FILE" env-required:"true"`
	// amount of workers
	Mappers int `env:"MAPPERS" env-required:"true"`
	// amount of reducers
	Reducers int `env:"REDUCERS" env-required:"true"`
	// transport to be used in the system
	Transport string `env:"TRANSPORT" env-required:"true"`
	// directory to map into master node
	MappedDir string `env:"MAPPED_DIRECTORY" env-required:"false" env-default:"data"`
	// time for all worker nodes to register on master
	RegisterDur time.Duration `env:"REGISTER_DURATION" env-required:"false" env-default:"30s"`

	// ================== MASTER ==================
	MasterHost            string        `env:"MASTER_HOST"                   env-required:"false" env-default:"localhost"`
	MasterPort            string        `env:"MASTER_PORT"                   env-required:"false" env-default:"8080"`
	MasterReadTimeout     time.Duration `env:"MASTER_READ_TIMEOUT"           env-required:"false" env-default:"3s"`
	MasterWriteTimeout    time.Duration `env:"MASTER_WRITE_TIMEOUT"          env-required:"false" env-default:"3s"`
	MasterShutdownTimeout time.Duration `env:"MASTER_SHUTDOWN_TIMEOUT"       env-required:"false" env-default:"0s"`
	// time for master to receive request from any worker
	MasterRequestAwaitDur time.Duration `env:"MASTER_REQUEST_AWAIT_DURATION" env-required:"false" env-default:"30s"`

	// ================== WORKER ==================
	WorkerHost            string        `env:"WORKER_HOST"              env-required:"false" env-default:"localhost"`
	WorkerReadTimeout     time.Duration `env:"WORKER_READ_TIMEOUT"      env-required:"false" env-default:"3s"`
	WorkerWriteTimeout    time.Duration `env:"WORKER_WRITE_TIMEOUT"     env-required:"false" env-default:"3s"`
	WorkerShutdownTimeout time.Duration `env:"WORKER_SHUTDOWN_TIMEOUT"  env-required:"false" env-default:"0s"`
	// time for worker node to set up
	WorkerSetupDur time.Duration `env:"WORKER_SETUP_DURATION" env-required:"false" env-default:"10s"`
	// time for worker to receive a new job
	WorkerAwaitDur time.Duration `env:"WORKER_AWAIT_DURATION" env-required:"false" env-default:"10s"`

	// ================== RABBITMQ ==================
	RabbitLogin        string `env:"RABBIT_LOGIN"         env-required:"false" env-default:"login"`
	RabbitPassword     string `env:"RABBIT_PASSWORD"      env-required:"false" env-defualt:"password"`
	RabbitHost         string `env:"RABBIT_HOST"          env-required:"false" env-default:"host"`
	RabbitPort         string `env:"RABBIT_PORT"          env-required:"false" env-default:"port"`
	RabbitMapperQueue  string `env:"RABBIT_MAPPER_QUEUE"  env-required:"false" env-default:"mapper"`
	RabbitReducerQueue string `env:"RABBIT_REDUCER_QUEUE" env-required:"false" env-default:"reducer"`
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

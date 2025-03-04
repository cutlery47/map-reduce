module github.com/cutlery47/map-reduce/master

go 1.23.2

require github.com/cutlery47/map-reduce/mapreduce v0.0.0

require (
	github.com/go-chi/chi/v5 v5.1.0
	github.com/rabbitmq/amqp091-go v1.10.0
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/ilyakaznacheev/cleanenv v1.5.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace github.com/cutlery47/map-reduce/mapreduce v0.0.0 => ../mapreduce

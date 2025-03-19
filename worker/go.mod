module github.com/cutlery47/map-reduce/worker

go 1.23.2

require (
	github.com/cutlery47/map-reduce/mapreduce v0.0.0
	github.com/go-chi/chi/v5 v5.1.0
	github.com/sirupsen/logrus v1.9.3
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/ilyakaznacheev/cleanenv v1.5.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace github.com/cutlery47/map-reduce/mapreduce v0.0.0 => ../mapreduce

module github.com/cutlery47/map-reduce/worker

go 1.23.2

require (
	github.com/cutlery47/map-reduce/mapreduce v0.0.0
	github.com/go-chi/chi/v5 v5.1.0
)

replace github.com/cutlery47/map-reduce/mapreduce v0.0.0 => ../mapreduce

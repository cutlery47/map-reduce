package mapreduce

var (
	MapperType  string = "MAPPER"
	ReducerType string = "REDUCER"
)

type WorkerRegisterRequest struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type WorkerRegisterResponse struct {
	Type string `json:"type"`
}

package mapreduce

// =================== DISTRIBUTION ========================

var (
	MapperDistr  string = "MAPPER"
	ReducerDistr string = "REDUCER"
	WorkerAck    string = "ACK"
)

type WorkerRegisterRequest struct {
	Host string
	Port string `json:"port"`
}

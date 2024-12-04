package mapreduce

// =================== DISTRIBUTION ========================

var (
	MapperDistr  string = "MAPPER"
	ReducerDistr string = "REDUCER"
	WorkerAck    string = "ACK"
)

type WorkerRegisterRequest struct {
	Port string `json:"port"`
}

package mapreduce

type Role string

type Addr struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

var (
	Mapper  Role = "MAPPER"
	Reducer Role = "REDUCER"
)

type RegisterRequest struct {
	Addr Addr `json:"addr"`
}

type RegisterResponse struct {
	Role Role `json:"role"`
}

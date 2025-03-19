package mapreduce

type Role string

var (
	Mapper  Role = "MAPPER"
	Reducer Role = "REDUCER"
)

type Code int

var (
	Success Code = 000001
	Failure Code = 000002
)

type Addr struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type RegisterRequest struct {
	Addr Addr `json:"addr"`
}

type RegisterResponse struct {
	Role Role `json:"role"`
}

type TerminateMessage struct {
	Code        Code   `json:"code"`
	Description string `json:"description"`
}

type TerminateRequest struct {
	Message TerminateMessage `json:"message"`
}

package requests

import (
	"github.com/cutlery47/map-reduce/mapreduce/models"
)

type RegisterRequest struct {
	Addr models.Addr `json:"addr"`
}

type RegisterResponse struct {
	Role models.Role `json:"role"`
}

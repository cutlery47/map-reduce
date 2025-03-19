package requests

import "github.com/cutlery47/map-reduce/mapreduce/models"

type TerminateRequest struct {
	Message models.TerminateMessage `json:"message"`
}

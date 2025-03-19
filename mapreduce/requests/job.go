package requests

import "github.com/cutlery47/map-reduce/mapreduce/models"

type JobRequest struct {
	Job models.Job `json:"job"`
}

type JobResponse struct {
	Result models.Result `json:"result"`
}

package models

import "github.com/google/uuid"

type Job struct {
	Id   uuid.UUID `json:"id"`
	Data []byte    `json:"data"`
}

type Result struct {
	JobId uuid.UUID `json:"job_id"`
	Code  Code      `json:"code"`
	Data  []byte    `json:"data"`
}

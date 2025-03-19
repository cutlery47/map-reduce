package models

type Role string

type Code int

type Addr struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type TerminateMessage struct {
	Code        Code   `json:"code"`
	Description string `json:"description"`
}

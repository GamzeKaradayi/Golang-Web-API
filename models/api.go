package models

type Api struct {
	Message  string `json:"message"`
	HasError bool   `json:"haserror"`
}

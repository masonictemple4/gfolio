package dtos

type ErrorReturn struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

package controllers

type FailureResponse struct {
	Error string `json:"error" example:"invalid input"`
}

type SuccessResponse struct {
	Message string      `json:"message" example:"OK"`
	Data    interface{} `json:"data,omitempty"`
}

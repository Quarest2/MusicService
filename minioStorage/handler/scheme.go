package handler

import (
	"minioStorage/client"
)

type Handler struct {
	minioService client.Client
}

type Services struct {
	minioService client.Client
}

type Handlers struct {
	minioHandler Handler // Пока у нас только один роут
}

type ObjectIdsDto struct {
	ObjectIDs []string `json:"objectIDs"`
}

type ErrorResponse struct {
	Error   string      `json:"error"`
	Status  int         `json:"code,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

type SuccessResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

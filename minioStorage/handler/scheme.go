package handler

import (
	"minioStorage/client"
)

type Handler struct {
	minioService client.Client
}

// Services структура всех сервисов, которые используются в хендлерах
// Это нужно чтобы мы могли использовать внутри хендлеров эти самые сервисы
type Services struct {
	minioService client.Client // Сервис у нас только один - minio, мы планируем его использовать, поэтому передаем
}

// Handlers структура всех хендлеров, которые используются для обозначения действия в роутах
type Handlers struct {
	minioHandler Handler // Пока у нас только один роут
}

// Нужен когда в body приходит много objectId - GetMany / DeleteMany
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

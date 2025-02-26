package client

import (
	"github.com/minio/minio-go/v7"
)

// Client интерфейс для взаимодействия с Minio
type Client interface {
	InitMinio() error                                     // Метод для инициализации подключения к Minio
	CreateOne(file FileDataType) (string, error)          // Метод для создания одного объекта в бакете Minio
	CreateMany(map[string]FileDataType) ([]string, error) // Метод для создания нескольких объектов в бакете Minio
	GetOne(objectID string) (string, error)               // Метод для получения одного объекта из бакета Minio
	GetMany(objectIDs []string) ([]string, error)         // Метод для получения нескольких объектов из бакета Minio
	DeleteOne(objectID string) error                      // Метод для удаления одного объекта из бакета Minio
	DeleteMany(objectIDs []string) error                  // Метод для удаления нескольких объектов из бакета Minio
}

// minioClient реализация интерфейса MinioClient
type minioClient struct {
	mc *minio.Client // Клиент Minio
}

type OperationError struct {
	ObjectID string
	Error    error
}

type FileDataType struct {
	FileName string
	Data     []byte
}

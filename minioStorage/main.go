package minioStorage

import (
	"github.com/gin-gonic/gin"
	"log"
	"minioStorage/client"
	"minioStorage/config"
	"minioStorage/handler"
)

func main() {
	config.LoadConfig()

	minioClient := client.NewMinioClient()
	err := minioClient.InitMinio()
	if err != nil {
		log.Fatalf("Ошибка инициализации Minio: %v", err)
	}

	_, s := handler.NewHandler(minioClient)

	router := gin.Default()

	s.RegisterRoutes(router)
}

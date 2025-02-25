package minioStorage

import (
	"github.com/gin-gonic/gin"
	"log"
	"minioStorage/client"
	"minioStorage/config"
	"minioStorage/handler"
)

func main() {
	// Загрузка конфигурации из файла .env
	config.LoadConfig()

	// Инициализация соединения с Minio
	minioClient := client.NewMinioClient()
	err := minioClient.InitMinio()
	if err != nil {
		log.Fatalf("Ошибка инициализации Minio: %v", err)
	}

	_, s := handler.NewHandler(minioClient)

	// Инициализация маршрутизатора Gin
	router := gin.Default()

	s.RegisterRoutes(router)

	// Запуск сервера Gin
	port := config.AppConfig.Port // Мы берем порт из конфига
	err = router.Run(":" + port)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера Gin: %v", err)
	}
}

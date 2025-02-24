package main

import (
	"MusicService/docs"
	"MusicService/minioStorage"
	"MusicService/restServer/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error with .env file: %v", err)
	}

	minioClient := minioStorage.NewMinioClient()
	err = minioClient.InitMinio()
	if err != nil {
		log.Fatalf("Error with initializing Minio: %v", err)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "OPTIONS"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Authorization"},
		AllowHeaders:     []string{"Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	docs.SwaggerInfo.Title = "Music Service API"
	docs.SwaggerInfo.Description = "This is an API of Music Service project."
	docs.SwaggerInfo.BasePath = "/api"

	swaggerUrl := ginSwagger.URL("./doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerUrl))

	apiGroup := r.Group("/api")

	apiGroup.GET("/ping", controllers.Ping)

	r.Run(":8080")
}

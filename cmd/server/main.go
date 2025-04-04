package main

import (
	"MusicService/internal/config"
	"MusicService/internal/controller"
	"MusicService/internal/middleware"
	"MusicService/internal/repository"
	"MusicService/internal/service"
	"MusicService/internal/storage"
	"MusicService/pkg/jwt"
	"MusicService/pkg/response"
	"log"

	_ "MusicService/docs" // Импорт сгенерированной документации
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Music Streaming Service API
// @version 1.0
// @description API для сервиса стриминга и селфхостинга музыки

// @contact.name API Support
// @contact.url http://example.com/support
// @contact.email support@music-service.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	minioClient, err := storage.NewMinioClient(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	err = minioClient.CreateBucket(cfg.MinIO.BucketName)
	if err != nil {
		log.Fatalf("Failed to create bucket: %v", err)
	}

	jwtService := jwt.NewJWTService(cfg.JWT.SecretKey)
	userRepo := repository.NewUserRepository(db)
	trackRepo := repository.NewTrackRepository(db)
	playlistRepo := repository.NewPlaylistRepository(db)

	authService := service.NewAuthService(userRepo, jwtService)
	userService := service.NewUserService(userRepo)
	trackService := service.NewTrackService(trackRepo, minioClient, cfg.MinIO.BucketName)
	playlistService := service.NewPlaylistService(playlistRepo, trackRepo)

	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)
	trackController := controller.NewTrackController(trackService)
	playlistController := controller.NewPlaylistController(playlistService)

	router := gin.Default()
	router.Use(response.CORSMiddleware())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(jwtService))
	{
		user := api.Group("/user")
		{
			user.GET("/profile", userController.GetProfile)
			user.PUT("/profile", userController.UpdateProfile)
		}

		track := api.Group("/tracks")
		{
			track.POST("", trackController.UploadTrack)
			track.GET("", trackController.GetAllTracks)
			track.GET("/user/:userId", trackController.GetUserTracks)
			track.GET("/:id", trackController.GetTrackByID)
			track.GET("/stream/:id", trackController.StreamTrack)
			track.DELETE("/:id", trackController.DeleteTrack)
			track.GET("/search", trackController.SearchTracks)
			track.GET("/:id/image", trackController.GetTrackImage)
		}

		playlist := api.Group("/playlists")
		{
			playlist.POST("", playlistController.CreatePlaylist)
			playlist.GET("", playlistController.GetUserPlaylists)
			playlist.GET("/:id", playlistController.GetPlaylistByID)
			playlist.PUT("/:id", playlistController.UpdatePlaylist)
			playlist.DELETE("/:id", playlistController.DeletePlaylist)
			playlist.POST("/:id/tracks", playlistController.AddTrackToPlaylist)
			playlist.DELETE("/:id/tracks/:trackId", playlistController.RemoveTrackFromPlaylist)
		}
	}

	log.Printf("Server is running on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

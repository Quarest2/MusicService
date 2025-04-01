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

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize MinIO client
	minioClient, err := storage.NewMinioClient(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	// Ensure bucket exists
	err = minioClient.CreateBucket(cfg.MinIO.BucketName)
	if err != nil {
		log.Fatalf("Failed to create bucket: %v", err)
	}

	// Initialize JWT service
	jwtService := jwt.NewJWTService(cfg.JWT.SecretKey)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	trackRepo := repository.NewTrackRepository(db)
	playlistRepo := repository.NewPlaylistRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, jwtService)
	userService := service.NewUserService(userRepo)
	trackService := service.NewTrackService(trackRepo, minioClient, cfg.MinIO.BucketName)
	playlistService := service.NewPlaylistService(playlistRepo, trackRepo)

	// Initialize controllers
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)
	trackController := controller.NewTrackController(trackService)
	playlistController := controller.NewPlaylistController(playlistService)

	// Setup router
	router := gin.Default()
	router.Use(response.CORSMiddleware())

	// Auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	// Authenticated routes
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(jwtService))
	{
		// User routes
		user := api.Group("/user")
		{
			user.GET("/profile", userController.GetProfile)
			user.PUT("/profile", userController.UpdateProfile)
		}

		// Track routes
		track := api.Group("/tracks")
		{
			track.POST("", trackController.UploadTrack)
			track.GET("", trackController.GetAllTracks)
			track.GET("/:id", trackController.GetTrackByID)
			track.GET("/stream/:id", trackController.StreamTrack)
			track.DELETE("/:id", trackController.DeleteTrack)
			track.GET("/search", trackController.SearchTracks)
		}

		// Playlist routes
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

	// Start server
	log.Printf("Server is running on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

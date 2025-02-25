package handler

import (
	"github.com/gin-gonic/gin"
	"minioStorage/client"
)

func NewMinioHandler(minioService client.Client) *Handler {
	return &Handler{
		minioService: minioService,
	}
}

// NewHandler создает экземпляр Handler с предоставленными сервисами
func NewHandler(minioService client.Client) (*Services, *Handlers) {
	return &Services{
			minioService: minioService,
		}, &Handlers{
			// инициируем Minio handler, который на вход получает minio service
			minioHandler: *NewMinioHandler(minioService),
		}
}

// RegisterRoutes - метод регистрации всех роутов в системе
func (h *Handlers) RegisterRoutes(router *gin.Engine) {

	// Здесь мы обозначили все эндпоинты системы с соответствующими хендлерами
	minioRoutes := router.Group("/files")
	{
		minioRoutes.POST("/", h.minioHandler.CreateOne)
		minioRoutes.POST("/many", h.minioHandler.CreateMany)

		minioRoutes.GET("/:objectID", h.minioHandler.GetOne)
		minioRoutes.GET("/many", h.minioHandler.GetMany)

		minioRoutes.DELETE("/:objectID", h.minioHandler.DeleteOne)
		minioRoutes.DELETE("/many", h.minioHandler.DeleteMany)
	}

}

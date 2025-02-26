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

func NewHandler(minioService client.Client) (*Services, *Handlers) {
	return &Services{
			minioService: minioService,
		}, &Handlers{
			minioHandler: *NewMinioHandler(minioService),
		}
}

func (h *Handlers) RegisterRoutes(router *gin.Engine) {

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

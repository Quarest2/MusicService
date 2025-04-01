package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response стандартная структура ответа API
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success отправляет успешный JSON-ответ
func Success(ctx *gin.Context, statusCode int, data interface{}) {
	ctx.JSON(statusCode, Response{
		Success: true,
		Data:    data,
	})
}

// Error отправляет JSON-ответ с ошибкой
func Error(ctx *gin.Context, statusCode int, errorMessage string) {
	ctx.JSON(statusCode, Response{
		Success: false,
		Error:   errorMessage,
	})
	ctx.Abort()
}

// CORSMiddleware middleware для обработки CORS заголовков
func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

// ValidationError отправляет ошибку валидации с деталями
func ValidationError(ctx *gin.Context, errors map[string]string) {
	ctx.JSON(http.StatusUnprocessableEntity, Response{
		Success: false,
		Error:   "Validation failed",
		Data:    errors,
	})
	ctx.Abort()
}

// PaginatedResponse структура для пагинированных ответов
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	TotalPages int         `json:"total_pages"`
}

// PaginatedSuccess отправляет успешный пагинированный ответ
func PaginatedSuccess(ctx *gin.Context, data interface{}, total int64, page int, perPage int) {
	totalPages := total / int64(perPage)
	if total%int64(perPage) > 0 {
		totalPages++
	}

	ctx.JSON(http.StatusOK, Response{
		Success: true,
		Data: PaginatedResponse{
			Data:       data,
			Total:      total,
			Page:       page,
			PerPage:    perPage,
			TotalPages: int(totalPages),
		},
	})
}

package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Success(ctx *gin.Context, statusCode int, data interface{}) {
	ctx.JSON(statusCode, Response{
		Success: true,
		Data:    data,
	})
}

func Error(ctx *gin.Context, statusCode int, errorMessage string) {
	ctx.JSON(statusCode, Response{
		Success: false,
		Error:   errorMessage,
	})
	ctx.Abort()
}

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

//func ValidationError(ctx *gin.Context, errors map[string]string) {
//	ctx.JSON(http.StatusUnprocessableEntity, Response{
//		Success: false,
//		Error:   "Validation failed",
//		Data:    errors,
//	})
//	ctx.Abort()
//}

//type PaginatedResponse struct {
//	Data       interface{} `json:"data"`
//	Total      int64       `json:"total"`
//	Page       int         `json:"page"`
//	PerPage    int         `json:"per_page"`
//	TotalPages int         `json:"total_pages"`
//}
//
//func PaginatedSuccess(ctx *gin.Context, data interface{}, total int64, page int, perPage int) {
//	totalPages := total / int64(perPage)
//	if total%int64(perPage) > 0 {
//		totalPages++
//	}
//
//	ctx.JSON(http.StatusOK, Response{
//		Success: true,
//		Data: PaginatedResponse{
//			Data:       data,
//			Total:      total,
//			Page:       page,
//			PerPage:    perPage,
//			TotalPages: int(totalPages),
//		},
//	})
//}

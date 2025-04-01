package middleware

import (
	"MusicService/pkg/jwt"
	"MusicService/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware проверяет JWT токен и добавляет userID в контекст
func AuthMiddleware(jwtService jwt.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Получаем заголовок Authorization
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(ctx, http.StatusUnauthorized, "Authorization header is required")
			ctx.Abort()
			return
		}

		// Проверяем формат заголовка (должен быть "Bearer <token>")
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			response.Error(ctx, http.StatusUnauthorized, "Invalid authorization header format")
			ctx.Abort()
			return
		}

		// Извлекаем токен
		token := headerParts[1]

		// Валидируем токен
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			response.Error(ctx, http.StatusUnauthorized, "Invalid token")
			ctx.Abort()
			return
		}

		// Добавляем userID в контекст
		ctx.Set("userID", claims.UserID)

		// Продолжаем выполнение
		ctx.Next()
	}
}

// AdminMiddleware проверяет, является ли пользователь администратором
// (В вашем ТЗ это не требуется, но добавлю для примера)
func AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// В реальном приложении здесь была бы проверка роли пользователя
		// Например, через запрос к БД или проверку claims.Role

		// Для примера просто возвращаем ошибку
		response.Error(ctx, http.StatusForbidden, "Admin access required")
		ctx.Abort()
	}
}

// CORSMiddleware (уже есть в response пакете, но можно дублировать здесь для удобства)
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

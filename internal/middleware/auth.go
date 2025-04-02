package middleware

import (
	"MusicService/pkg/jwt"
	"MusicService/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService jwt.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(ctx, http.StatusUnauthorized, "Authorization header is required")
			ctx.Abort()
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			response.Error(ctx, http.StatusUnauthorized, "Invalid authorization header format")
			ctx.Abort()
			return
		}

		token := headerParts[1]

		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			response.Error(ctx, http.StatusUnauthorized, "Invalid token")
			ctx.Abort()
			return
		}

		ctx.Set("userID", claims.UserID)

		ctx.Next()
	}
}

//
//func AdminMiddleware() gin.HandlerFunc {
//	return func(ctx *gin.Context) {
//
//		response.Error(ctx, http.StatusForbidden, "Admin access required")
//		ctx.Abort()
//	}
//}

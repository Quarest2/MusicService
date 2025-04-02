package controller

import (
	"MusicService/internal/model"
	"MusicService/internal/service"
	"MusicService/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Register godoc
// @Summary Регистрация пользователя
// @Description Создает нового пользователя
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.RegisterRequest true "Данные регистрации"
// @Success 201 {object} model.UserResponse
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req model.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := c.authService.Register(&req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "user already exists" {
			status = http.StatusConflict
		}
		response.Error(ctx, status, err.Error())
		return
	}

	response.Success(ctx, http.StatusCreated, user)
}

// Login godoc
// @Summary Авторизация пользователя
// @Description Вход в систему и получение JWT токена
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Данные авторизации"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req model.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request payload")
		return
	}

	token, err := c.authService.Login(&req)
	if err != nil {
		response.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, gin.H{"token": token})
}

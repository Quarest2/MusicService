package controller

import (
	"MusicService/internal/model"
	"MusicService/internal/service"
	"MusicService/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService: userService}
}

// GetProfile godoc
// @Summary Получить профиль пользователя
// @Description Возвращает информацию о текущем авторизованном пользователе
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.UserResponse
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/user/profile [get]
func (c *UserController) GetProfile(ctx *gin.Context) {
	userID := ctx.GetUint("userID")

	profile, err := c.userService.GetProfile(userID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		response.Error(ctx, status, "Failed to get profile")
		return
	}

	response.Success(ctx, http.StatusOK, profile)
}

// UpdateProfile godoc
// @Summary Обновить профиль пользователя
// @Description Обновляет информацию о текущем авторизованном пользователе
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.UserResponse true "Данные для обновления"
// @Success 200 {object} model.UserResponse
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/user/profile [put]
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID := ctx.GetUint("userID")

	var update model.UserResponse
	if err := ctx.ShouldBindJSON(&update); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request payload")
		return
	}

	profile, err := c.userService.UpdateProfile(userID, &update)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "email already exists":
			status = http.StatusConflict
		case "username already exists":
			status = http.StatusConflict
		}
		response.Error(ctx, status, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, profile)
}

//// ChangePassword godoc
//// @Summary Изменить пароль
//// @Description Изменяет пароль текущего авторизованного пользователя
//// @Tags User
//// @Accept json
//// @Produce json
//// @Security BearerAuth
//// @Param request body model.ChangePasswordRequest true "Данные для смены пароля"
//// @Success 200 {object} response.Response
//// @Failure 400 {object} response.Response
//// @Failure 401 {object} response.Response
//// @Failure 403 {object} response.Response
//// @Failure 500 {object} response.Response
//// @Router /api/user/change-password [post]
//func (c *UserController) ChangePassword(ctx *gin.Context) {
//	userID := ctx.GetUint("userID")
//
//	var req model.ChangePasswordRequest
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		response.Error(ctx, http.StatusBadRequest, "Invalid request payload")
//		return
//	}
//
//	if err := c.userService.ChangePassword(userID, &req); err != nil {
//		status := http.StatusInternalServerError
//		if err.Error() == "current password is incorrect" {
//			status = http.StatusForbidden
//		}
//		response.Error(ctx, status, err.Error())
//		return
//	}
//
//	response.Success(ctx, http.StatusOK, gin.H{"message": "Password changed successfully"})
//}

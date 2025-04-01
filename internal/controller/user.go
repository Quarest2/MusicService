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

func (c *UserController) GetProfile(ctx *gin.Context) {
	userID := ctx.GetUint("userID")

	profile, err := c.userService.GetProfile(userID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to get profile")
		return
	}

	response.Success(ctx, http.StatusOK, profile)
}

func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID := ctx.GetUint("userID")

	var update model.UserResponse
	if err := ctx.ShouldBindJSON(&update); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request payload")
		return
	}

	profile, err := c.userService.UpdateProfile(userID, &update)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	response.Success(ctx, http.StatusOK, profile)
}

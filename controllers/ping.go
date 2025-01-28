package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary Ping server
// @Description ping server
// @ID ping-server
// @Tags api
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} FailureResponse
// @Router /ping [get]
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, SuccessResponse{Message: "pong"})
}

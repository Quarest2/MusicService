package controller

import (
	"MusicService/internal/model"
	"MusicService/internal/service"
	"MusicService/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PlaylistController struct {
	playlistService service.PlaylistService
}

func NewPlaylistController(playlistService service.PlaylistService) *PlaylistController {
	return &PlaylistController{playlistService: playlistService}
}

func (c *PlaylistController) CreatePlaylist(ctx *gin.Context) {
	userID := ctx.GetUint("userID")

	var req model.PlaylistRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request payload")
		return
	}

	playlist, err := c.playlistService.CreatePlaylist(&req, userID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to create playlist")
		return
	}

	response.Success(ctx, http.StatusCreated, playlist)
}

func (c *PlaylistController) GetUserPlaylists(ctx *gin.Context) {
	userID := ctx.GetUint("userID")

	playlists, err := c.playlistService.GetUserPlaylists(userID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to get playlists")
		return
	}

	response.Success(ctx, http.StatusOK, playlists)
}

func (c *PlaylistController) GetPlaylistByID(ctx *gin.Context) {
	var uri struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := ctx.ShouldBindUri(&uri); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid playlist ID")
		return
	}

	playlist, err := c.playlistService.GetPlaylistByID(uri.ID)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "Playlist not found")
		return
	}

	response.Success(ctx, http.StatusOK, playlist)
}

func (c *PlaylistController) UpdatePlaylist(ctx *gin.Context) {
	var uri struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := ctx.ShouldBindUri(&uri); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid playlist ID")
		return
	}

	var req model.PlaylistRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updatedPlaylist, err := c.playlistService.UpdatePlaylist(uri.ID, &req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to update playlist")
		return
	}

	response.Success(ctx, http.StatusOK, updatedPlaylist)
}

func (c *PlaylistController) DeletePlaylist(ctx *gin.Context) {
	var uri struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := ctx.ShouldBindUri(&uri); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid playlist ID")
		return
	}

	if err := c.playlistService.DeletePlaylist(uri.ID); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to delete playlist")
		return
	}

	response.Success(ctx, http.StatusOK, gin.H{"message": "Playlist deleted successfully"})
}

func (c *PlaylistController) AddTrackToPlaylist(ctx *gin.Context) {
	var uri struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := ctx.ShouldBindUri(&uri); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid playlist ID")
		return
	}

	var req model.AddTrackToPlaylistRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := c.playlistService.AddTrackToPlaylist(uri.ID, &req); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to add track to playlist")
		return
	}

	response.Success(ctx, http.StatusOK, gin.H{"message": "Track added to playlist successfully"})
}

func (c *PlaylistController) RemoveTrackFromPlaylist(ctx *gin.Context) {
	var uri struct {
		ID      uint `uri:"id" binding:"required"`
		TrackID uint `uri:"trackId" binding:"required"`
	}

	if err := ctx.ShouldBindUri(&uri); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request parameters")
		return
	}

	if err := c.playlistService.RemoveTrackFromPlaylist(uri.ID, uri.TrackID); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to remove track from playlist")
		return
	}

	response.Success(ctx, http.StatusOK, gin.H{"message": "Track removed from playlist successfully"})
}

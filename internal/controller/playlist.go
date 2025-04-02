package controller

import (
	"MusicService/internal/model"
	"MusicService/internal/service"
	"MusicService/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PlaylistController struct {
	playlistService service.PlaylistService
}

func NewPlaylistController(playlistService service.PlaylistService) *PlaylistController {
	return &PlaylistController{
		playlistService: playlistService,
	}
}

// CreatePlaylist godoc
// @Summary Создать плейлист
// @Description Создает новый плейлист для текущего пользователя
// @Tags Playlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.PlaylistRequest true "Данные плейлиста"
// @Success 201 {object} model.PlaylistResponse
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/playlists [post]
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

// GetUserPlaylists godoc
// @Summary Получить плейлисты пользователя
// @Description Возвращает все плейлисты текущего пользователя
// @Tags Playlists
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.PlaylistResponse
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/playlists [get]
func (c *PlaylistController) GetUserPlaylists(ctx *gin.Context) {
	userID := ctx.GetUint("userID")

	playlists, err := c.playlistService.GetUserPlaylists(userID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to get playlists")
		return
	}

	response.Success(ctx, http.StatusOK, playlists)
}

// GetPlaylistByID godoc
// @Summary Получить плейлист по ID
// @Description Возвращает плейлист с указанным ID
// @Tags Playlists
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID плейлиста"
// @Success 200 {object} model.PlaylistResponse
// @Failure 400 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/playlists/{id} [get]
func (c *PlaylistController) GetPlaylistByID(ctx *gin.Context) {
	playlistID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid playlist ID")
		return
	}

	playlist, err := c.playlistService.GetPlaylistByID(uint(playlistID))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "playlist not found" {
			status = http.StatusNotFound
		}
		response.Error(ctx, status, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, playlist)
}

// UpdatePlaylist godoc
// @Summary Обновить плейлист
// @Description Обновляет информацию о плейлисте
// @Tags Playlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID плейлиста"
// @Param request body model.PlaylistRequest true "Новые данные плейлиста"
// @Success 200 {object} model.PlaylistResponse
// @Failure 400 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/playlists/{id} [put]
func (c *PlaylistController) UpdatePlaylist(ctx *gin.Context) {
	playlistID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid playlist ID")
		return
	}

	var req model.PlaylistRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updatedPlaylist, err := c.playlistService.UpdatePlaylist(uint(playlistID), &req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to update playlist")
		return
	}

	response.Success(ctx, http.StatusOK, updatedPlaylist)
}

// DeletePlaylist godoc
// @Summary Удалить плейлист
// @Description Удаляет плейлист с указанным ID
// @Tags Playlists
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID плейлиста"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/playlists/{id} [delete]
func (c *PlaylistController) DeletePlaylist(ctx *gin.Context) {
	playlistID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid playlist ID")
		return
	}

	if err := c.playlistService.DeletePlaylist(uint(playlistID)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to delete playlist")
		return
	}

	response.Success(ctx, http.StatusOK, gin.H{"message": "Playlist deleted successfully"})
}

// AddTrackToPlaylist godoc
// @Summary Добавить трек в плейлист
// @Description Добавляет трек в указанный плейлист
// @Tags Playlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID плейлиста"
// @Param request body model.AddTrackToPlaylistRequest true "ID трека"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/playlists/{id}/tracks [post]
func (c *PlaylistController) AddTrackToPlaylist(ctx *gin.Context) {
	playlistID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid playlist ID")
		return
	}

	var req model.AddTrackToPlaylistRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := c.playlistService.AddTrackToPlaylist(uint(playlistID), &req); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "track not found" {
			status = http.StatusNotFound
		}
		response.Error(ctx, status, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, gin.H{"message": "Track added to playlist successfully"})
}

// RemoveTrackFromPlaylist godoc
// @Summary Удалить трек из плейлиста
// @Description Удаляет трек из указанного плейлиста
// @Tags Playlists
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID плейлиста"
// @Param trackId path int true "ID трека"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/playlists/{id}/tracks/{trackId} [delete]
func (c *PlaylistController) RemoveTrackFromPlaylist(ctx *gin.Context) {
	playlistID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request parameters")
		return
	}

	trackID, err := strconv.ParseUint(ctx.Param("trackId"), 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Failed to remove track from playlist")
		return
	}

	if err := c.playlistService.RemoveTrackFromPlaylist(uint(playlistID), uint(trackID)); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "track not found in playlist" {
			status = http.StatusNotFound
		}
		response.Error(ctx, status, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, gin.H{"message": "Track removed from playlist successfully"})
}

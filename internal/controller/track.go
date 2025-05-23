package controller

import (
	"MusicService/internal/model"
	"MusicService/internal/service"
	"MusicService/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TrackController struct {
	trackService service.TrackService
	statsService service.StatsService
}

func NewTrackController(trackService service.TrackService, statsService service.StatsService) *TrackController {
	return &TrackController{trackService: trackService, statsService: statsService}
}

// UploadTrack загружает новый трек
// @Summary Загрузить новый трек
// @Description Загружает аудиофайл и создает запись о треке
// @Tags Tracks
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Аудиофайл"
// @Param image formData file false "Изображение"
// @Param title formData string true "Название трека"
// @Param artist formData string true "Исполнитель"
// @Param album formData string false "Альбом"
// @Param genre formData string false "Жанр"
// @Success 201 {object} model.TrackResponse
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/tracks [post]
func (c *TrackController) UploadTrack(ctx *gin.Context) {
	userID := ctx.GetUint("userID")

	audioFile, err := ctx.FormFile("file")
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "File is required")
		return
	}

	imageFile, err := ctx.FormFile("image")

	var req model.TrackUploadRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid request payload")
		return
	}

	track, err := c.trackService.UploadTrack(audioFile, imageFile, &req, userID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to upload track")
		return
	}

	response.Success(ctx, http.StatusCreated, track)
}

// GetAllTracks возвращает все треки
// @Summary Получить все треки
// @Description Возвращает список всех треков в системе
// @Tags Tracks
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.TrackResponse
// @Failure 500 {object} response.Response
// @Router /api/tracks [get]
func (c *TrackController) GetAllTracks(ctx *gin.Context) {
	tracks, err := c.trackService.GetAllTracks()
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to get tracks")
		return
	}

	response.Success(ctx, http.StatusOK, tracks)
}

// GetTrackByID возвращает трек по ID
// @Summary Получить трек по ID
// @Description Возвращает информацию о конкретном треке
// @Tags Tracks
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID трека"
// @Success 200 {object} model.TrackResponse
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/tracks/{id} [get]
func (c *TrackController) GetTrackByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid track ID")
		return
	}

	track, err := c.trackService.GetTrackByID(uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "Track not found")
		return
	}

	response.Success(ctx, http.StatusOK, track)
}

// StreamTrack возвращает аудиопоток трека
// @Summary Воспроизвести трек
// @Description Возвращает аудиопоток для проигрывания трека
// @Tags Tracks
// @Produce audio/mpeg
// @Security BearerAuth
// @Param id path int true "ID трека"
// @Success 200 {file} binary
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/tracks/stream/{id} [get]
func (c *TrackController) StreamTrack(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid track ID")
		return
	}

	reader, contentType, err := c.trackService.StreamTrack(uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "Track not found")
		return
	}
	defer reader.Close()

	if err = c.statsService.RecordTrackPlay(userID, uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to record play of the track")
	}

	ctx.DataFromReader(http.StatusOK, -1, contentType, reader, nil)
}

// DeleteTrack удаляет трек
// @Summary Удалить трек
// @Description Удаляет трек по ID (только для владельца)
// @Tags Tracks
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID трека"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/tracks/{id} [delete]
func (c *TrackController) DeleteTrack(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid track ID")
		return
	}

	if err := c.trackService.DeleteTrack(uint(id)); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to delete track")
		return
	}

	response.Success(ctx, http.StatusOK, gin.H{"message": "Track deleted successfully"})
}

// SearchTracks ищет треки по параметрам
// @Summary Поиск треков
// @Description Поиск треков по названию, исполнителю, альбому или жанру
// @Tags Tracks
// @Produce json
// @Security BearerAuth
// @Param q query string false "Поисковый запрос"
// @Param artist query string false "Исполнитель"
// @Param album query string false "Альбом"
// @Param genre query string false "Жанр"
// @Success 200 {array} model.TrackResponse
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/tracks/search [get]
func (c *TrackController) SearchTracks(ctx *gin.Context) {
	var params model.TrackSearchParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid search parameters")
		return
	}

	tracks, err := c.trackService.SearchTracks(params)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to search tracks")
		return
	}

	response.Success(ctx, http.StatusOK, tracks)
}

// GetTrackImage godoc
// @Summary Получить изображение трека
// @Tags Tracks
// @Produce image/*
// @Security BearerAuth
// @Param id path int true "ID трека"
// @Success 200 {file} byte
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/tracks/{id}/image [get]
func (c *TrackController) GetTrackImage(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid track ID")
		return
	}

	reader, contentType, err := c.trackService.GetTrackImage(uint(id))
	if err != nil {
		response.Error(ctx, http.StatusNotFound, "Image not found")
		return
	}
	defer reader.Close()

	ctx.DataFromReader(http.StatusOK, -1, contentType, reader, nil)
}

// GetUserTracks возвращает все треки пользователя
// @Summary Получить все треки пользователя
// @Description Возвращает список всех треков пользователя в системе
// @Tags Tracks
// @Produce json
// @Security BearerAuth
// @Param userId path int true "ID пользователя"
// @Success 200 {array} model.TrackResponse
// @Failure 500 {object} response.Response
// @Router /api/tracks/user/{userId} [get]
func (c *TrackController) GetUserTracks(ctx *gin.Context) {
	userId, err := strconv.ParseUint(ctx.Param("userId"), 10, 32)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid user ID")
		return
	}

	tracks, err := c.trackService.GetUserTracks(uint(userId))
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to get tracks")
		return
	}

	response.Success(ctx, http.StatusOK, tracks)
}

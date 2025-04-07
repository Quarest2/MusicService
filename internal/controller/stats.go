package controller

import (
	"MusicService/internal/service"
	"MusicService/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type StatsController struct {
	statsService service.StatsService
}

func NewStatsController(statsService service.StatsService) *StatsController {
	return &StatsController{
		statsService: statsService,
	}
}

// GetTrackPlaysCount godoc
// @Summary Получить кол-во прослушиваний у треков
// @Description Получить кол-во прослушиваний у треков для авторизированного пользователя
// @Tags Stats
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} model.TrackPlayStats
// @Failure 500 {object} response.Response
// @Router /api/stats/track-plays [get]
func (c *StatsController) GetTrackPlaysCount(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)

	stats, err := c.statsService.GetTrackPlaysStats(userID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to get stats")
		return
	}

	response.Success(ctx, http.StatusOK, stats)
}

// GetArtistPlaysCount godoc
// @Summary Получить кол-во прослушиваний у исполнителей
// @Description Получить кол-во прослушиваний у исполнителей для авторизированного пользователя
// @Tags Stats
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} model.ArtistPlayStats
// @Failure 500 {object} response.Response
// @Router /api/stats/artist-plays [get]
func (c *StatsController) GetArtistPlaysCount(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)

	stats, err := c.statsService.GetArtistPlaysStats(userID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to get stats")
		return
	}

	response.Success(ctx, http.StatusOK, stats)
}

// GetRecentTracks godoc
// @Summary Получить последние прослушанные треки
// @Description Получить последние (50) прослушанные треки
// @Tags Stats
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param limit query int false "Number of tracks to return (default 50)"
// @Success 200 {array} model.TrackResponse
// @Failure 500 {object} response.Response
// @Router /api/stats/recent-tracks [get]
func (c *StatsController) GetRecentTracks(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "50"))

	tracks, err := c.statsService.GetRecentTracks(userID, limit)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to get recent tracks")
		return
	}

	response.Success(ctx, http.StatusOK, tracks)
}

// GetRecentArtists godoc
// @Summary Получить последних прослушанных исполнителей
// @Description Получить последних (5) прослушанных исполнителей
// @Tags Stats
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param limit query int false "Number of artists to return (default 5)"
// @Success 200 {array} string
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/stats/recent-artists [get]
func (c *StatsController) GetRecentArtists(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "5"))

	artists, err := c.statsService.GetRecentArtists(userID, limit)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to get recent artists")
		return
	}

	response.Success(ctx, http.StatusOK, artists)
}

package repository

import (
	"MusicService/internal/model"
	"gorm.io/gorm"
)

type StatsRepository interface {
	GetTrackPlaysStats(userID uint) ([]model.TrackPlayStats, error)
	GetArtistPlaysStats(userID uint) ([]model.ArtistPlayStats, error)
	GetRecentTracks(userID uint, limit int) ([]model.Track, error)
	GetRecentArtists(userID uint, limit int) ([]string, error)
	CreateListeningHistory(history *model.ListeningHistory) error
}

type statsRepository struct {
	db *gorm.DB
}

func NewStatsRepository(db *gorm.DB) StatsRepository {
	return &statsRepository{db: db}
}

func (r *statsRepository) GetTrackPlaysStats(userID uint) ([]model.TrackPlayStats, error) {
	var stats []model.TrackPlayStats

	err := r.db.Model(&model.ListeningHistory{}).
		Select("tracks.id as track_id, tracks.title, tracks.artist, count(listening_histories.id) as play_count").
		Joins("left join tracks on tracks.id = listening_histories.track_id").
		Where("listening_histories.user_id = ?", userID).
		Group("tracks.id, tracks.title, tracks.artist").
		Order("play_count desc").
		Scan(&stats).Error

	return stats, err
}

func (r *statsRepository) GetArtistPlaysStats(userID uint) ([]model.ArtistPlayStats, error) {
	var stats []model.ArtistPlayStats

	err := r.db.Model(&model.ListeningHistory{}).
		Select("tracks.artist, count(listening_histories.id) as play_count").
		Joins("join tracks on tracks.id = listening_histories.track_id").
		Where("listening_histories.user_id = ?", userID).
		Group("tracks.artist").
		Order("play_count desc").
		Scan(&stats).Error

	return stats, err
}

func (r *statsRepository) GetRecentTracks(userID uint, limit int) ([]model.Track, error) {
	var tracks []model.Track

	err := r.db.Joins("JOIN listening_histories on listening_histories.track_id = tracks.id").
		Where("listening_histories.user_id = ?", userID).
		Order("listening_histories.created_at desc").
		Limit(limit).
		Find(&tracks).Error

	return tracks, err
}

func (r *statsRepository) GetRecentArtists(userID uint, limit int) ([]string, error) {
	var artists []string

	err := r.db.Model(&model.ListeningHistory{}).
		Select("distinct tracks.artist").
		Joins("join tracks on tracks.id = listening_histories.track_id").
		Where("listening_histories.user_id = ?", userID).
		Group("tracks.artist").
		Order("max(listening_histories.created_at) desc").
		Limit(limit).
		Pluck("tracks.artist", &artists).Error

	return artists, err
}

func (r *statsRepository) CreateListeningHistory(history *model.ListeningHistory) error {
	return r.db.Create(history).Error
}

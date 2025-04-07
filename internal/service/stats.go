package service

import (
	"MusicService/internal/model"
	"MusicService/internal/repository"
)

type StatsService interface {
	GetTrackPlaysStats(userID uint) ([]model.TrackPlayStats, error)
	GetArtistPlaysStats(userID uint) ([]model.ArtistPlayStats, error)
	GetRecentTracks(id uint, limit int) ([]model.Track, error)
	GetRecentArtists(id uint, limit int) ([]string, error)
	RecordTrackPlay(userID, trackID uint) error
}

type statsService struct {
	statsRepo repository.StatsRepository
}

func NewStatsService(statsRepo repository.StatsRepository) StatsService {
	return &statsService{
		statsRepo: statsRepo,
	}
}

func (s *statsService) GetTrackPlaysStats(userID uint) ([]model.TrackPlayStats, error) {
	return s.statsRepo.GetTrackPlaysStats(userID)
}

func (s *statsService) GetArtistPlaysStats(userID uint) ([]model.ArtistPlayStats, error) {
	return s.statsRepo.GetArtistPlaysStats(userID)
}

func (s *statsService) GetRecentTracks(userID uint, limit int) ([]model.Track, error) {
	return s.statsRepo.GetRecentTracks(userID, limit)
}

func (s *statsService) GetRecentArtists(userID uint, limit int) ([]string, error) {
	return s.statsRepo.GetRecentArtists(userID, limit)
}

func (s *statsService) RecordTrackPlay(userID, trackID uint) error {
	history := model.ListeningHistory{
		UserID:  userID,
		TrackID: trackID,
	}

	return s.statsRepo.CreateListeningHistory(&history)
}

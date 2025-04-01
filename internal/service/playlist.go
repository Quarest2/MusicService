package service

import (
	"MusicService/internal/model"
	"MusicService/internal/repository"
	"time"
)

type PlaylistService interface {
	CreatePlaylist(req *model.PlaylistRequest, userID uint) (*model.PlaylistResponse, error)
	GetUserPlaylists(userID uint) ([]model.PlaylistResponse, error)
	GetPlaylistByID(id uint) (*model.PlaylistResponse, error)
	UpdatePlaylist(id uint, req *model.PlaylistRequest) (*model.PlaylistResponse, error)
	DeletePlaylist(id uint) error
	AddTrackToPlaylist(playlistID uint, req *model.AddTrackToPlaylistRequest) error
	RemoveTrackFromPlaylist(playlistID uint, trackID uint) error
}

type playlistService struct {
	playlistRepo repository.PlaylistRepository
	trackRepo    repository.TrackRepository
}

func NewPlaylistService(playlistRepo repository.PlaylistRepository, trackRepo repository.TrackRepository) PlaylistService {
	return &playlistService{
		playlistRepo: playlistRepo,
		trackRepo:    trackRepo,
	}
}

func (s *playlistService) CreatePlaylist(req *model.PlaylistRequest, userID uint) (*model.PlaylistResponse, error) {
	playlist := &model.Playlist{
		Name:        req.Name,
		Description: req.Description,
		UserID:      userID,
	}

	if err := s.playlistRepo.Create(playlist); err != nil {
		return nil, err
	}

	return &model.PlaylistResponse{
		ID:          playlist.ID,
		Name:        playlist.Name,
		Description: playlist.Description,
		Tracks:      []model.TrackResponse{},
		CreatedAt:   playlist.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *playlistService) GetUserPlaylists(userID uint) ([]model.PlaylistResponse, error) {
	playlists, err := s.playlistRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var response []model.PlaylistResponse
	for _, playlist := range playlists {
		tracks, err := s.trackRepo.GetByPlaylistID(playlist.ID)
		if err != nil {
			return nil, err
		}

		var trackResponses []model.TrackResponse
		for _, track := range tracks {
			trackResponses = append(trackResponses, model.TrackResponse{
				ID:        track.ID,
				Title:     track.Title,
				Artist:    track.Artist,
				Album:     track.Album,
				Genre:     track.Genre,
				Duration:  track.Duration,
				CreatedAt: track.CreatedAt.Format(time.RFC3339),
			})
		}

		response = append(response, model.PlaylistResponse{
			ID:          playlist.ID,
			Name:        playlist.Name,
			Description: playlist.Description,
			Tracks:      trackResponses,
			CreatedAt:   playlist.CreatedAt.Format(time.RFC3339),
		})
	}

	return response, nil
}

func (s *playlistService) GetPlaylistByID(id uint) (*model.PlaylistResponse, error) {
	playlist, err := s.playlistRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	tracks, err := s.trackRepo.GetByPlaylistID(playlist.ID)
	if err != nil {
		return nil, err
	}

	var trackResponses []model.TrackResponse
	for _, track := range tracks {
		trackResponses = append(trackResponses, model.TrackResponse{
			ID:        track.ID,
			Title:     track.Title,
			Artist:    track.Artist,
			Album:     track.Album,
			Genre:     track.Genre,
			Duration:  track.Duration,
			CreatedAt: track.CreatedAt.Format(time.RFC3339),
		})
	}

	return &model.PlaylistResponse{
		ID:          playlist.ID,
		Name:        playlist.Name,
		Description: playlist.Description,
		Tracks:      trackResponses,
		CreatedAt:   playlist.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *playlistService) UpdatePlaylist(id uint, req *model.PlaylistRequest) (*model.PlaylistResponse, error) {
	playlist, err := s.playlistRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	playlist.Name = req.Name
	playlist.Description = req.Description

	if err := s.playlistRepo.Update(playlist); err != nil {
		return nil, err
	}

	return s.GetPlaylistByID(id)
}

func (s *playlistService) DeletePlaylist(id uint) error {
	return s.playlistRepo.Delete(id)
}

func (s *playlistService) AddTrackToPlaylist(playlistID uint, req *model.AddTrackToPlaylistRequest) error {
	return s.playlistRepo.AddTrack(playlistID, req.TrackID)
}

func (s *playlistService) RemoveTrackFromPlaylist(playlistID uint, trackID uint) error {
	return s.playlistRepo.RemoveTrack(playlistID, trackID)
}

package repository

import (
	"MusicService/internal/model"

	"gorm.io/gorm"
)

type PlaylistRepository interface {
	Create(playlist *model.Playlist) error
	GetByID(id uint) (*model.Playlist, error)
	GetByUserID(userID uint) ([]model.Playlist, error)
	Update(playlist *model.Playlist) error
	Delete(id uint) error
	AddTrack(playlistID uint, trackID uint) error
	RemoveTrack(playlistID uint, trackID uint) error
}

type playlistRepository struct {
	db *gorm.DB
}

func NewPlaylistRepository(db *gorm.DB) PlaylistRepository {
	return &playlistRepository{db: db}
}

func (r *playlistRepository) Create(playlist *model.Playlist) error {
	return r.db.Create(playlist).Error
}

func (r *playlistRepository) GetByID(id uint) (*model.Playlist, error) {
	var playlist model.Playlist
	err := r.db.Preload("Tracks").First(&playlist, id).Error
	return &playlist, err
}

func (r *playlistRepository) GetByUserID(userID uint) ([]model.Playlist, error) {
	var playlists []model.Playlist
	err := r.db.Where("user_id = ?", userID).Preload("Tracks").Find(&playlists).Error
	return playlists, err
}

func (r *playlistRepository) Update(playlist *model.Playlist) error {
	return r.db.Save(playlist).Error
}

func (r *playlistRepository) Delete(id uint) error {
	return r.db.Delete(&model.Playlist{}, id).Error
}

func (r *playlistRepository) AddTrack(playlistID uint, trackID uint) error {
	playlist, err := r.GetByID(playlistID)
	if err != nil {
		return err
	}

	track := &model.Track{Model: gorm.Model{ID: trackID}}
	return r.db.Model(playlist).Association("Tracks").Append(track)
}

func (r *playlistRepository) RemoveTrack(playlistID uint, trackID uint) error {
	playlist, err := r.GetByID(playlistID)
	if err != nil {
		return err
	}

	track := &model.Track{Model: gorm.Model{ID: trackID}}
	return r.db.Model(playlist).Association("Tracks").Delete(track)
}

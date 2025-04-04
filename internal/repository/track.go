package repository

import (
	"MusicService/internal/model"
	"gorm.io/gorm"
)

type TrackRepository interface {
	Create(track *model.Track) error
	GetByID(id uint) (*model.Track, error)
	GetAll() ([]model.Track, error)
	GetUserTracks(userId uint) ([]model.Track, error)
	Delete(id uint) error
	Search(params model.TrackSearchParams) ([]model.Track, error)
	GetByPlaylistID(playlistID uint) ([]model.Track, error)
}

type trackRepository struct {
	db *gorm.DB
}

func NewTrackRepository(db *gorm.DB) TrackRepository {
	return &trackRepository{db: db}
}

func (r *trackRepository) Create(track *model.Track) error {
	return r.db.Create(track).Error
}

func (r *trackRepository) GetByID(id uint) (*model.Track, error) {
	var track model.Track
	err := r.db.First(&track, id).Error
	return &track, err
}

func (r *trackRepository) GetAll() ([]model.Track, error) {
	var tracks []model.Track
	err := r.db.Find(&tracks).Error
	return tracks, err
}

func (r *trackRepository) GetUserTracks(uploadedBy uint) ([]model.Track, error) {
	var tracks []model.Track
	query := r.db.Model(&model.Track{})

	query = query.Where("uploaded_by = ?", uploadedBy)

	err := query.Find(&tracks).Error
	return tracks, err
}

func (r *trackRepository) Delete(id uint) error {
	return r.db.Delete(&model.Track{}, id).Error
}

func (r *trackRepository) Search(params model.TrackSearchParams) ([]model.Track, error) {
	var tracks []model.Track
	query := r.db.Model(&model.Track{})

	if params.Query != "" {
		query = query.Where("title LIKE ? OR artist LIKE ?", "%"+params.Query+"%", "%"+params.Query+"%")
	}
	if params.Artist != "" {
		query = query.Where("artist = ?", params.Artist)
	}
	if params.Album != "" {
		query = query.Where("album = ?", params.Album)
	}
	if params.Genre != "" {
		query = query.Where("genre = ?", params.Genre)
	}

	err := query.Find(&tracks).Error
	return tracks, err
}

func (r *trackRepository) GetByPlaylistID(playlistID uint) ([]model.Track, error) {
	var tracks []model.Track
	err := r.db.Model(&model.Playlist{}).Where("id = ?", playlistID).Association("Tracks").Find(&tracks)
	return tracks, err
}

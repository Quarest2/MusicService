package model

import "gorm.io/gorm"

type Playlist struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	UserID      uint    `gorm:"not null"`
	Tracks      []Track `gorm:"many2many:playlist_tracks;"`
}

type PlaylistRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type PlaylistResponse struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Tracks      []TrackResponse `json:"tracks"`
	CreatedAt   string          `json:"createdAt"`
}

type AddTrackToPlaylistRequest struct {
	TrackID uint `json:"trackId" binding:"required"`
}

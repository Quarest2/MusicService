package model

import (
	"gorm.io/gorm"
)

type ListeningHistory struct {
	gorm.Model
	UserID  uint  `json:"user_id"`
	TrackID uint  `json:"track_id"`
	Track   Track `json:"-" gorm:"foreignKey:TrackID"`
}

type TrackPlayStats struct {
	ID        uint   `json:"id" db:"id"`
	Title     string `json:"title" db:"title"`
	Artist    string `json:"artist" db:"artist"`
	PlayCount int    `json:"play_count" db:"play_count"`
}

type ArtistPlayStats struct {
	Artist    string `json:"artist" db:"artist"`
	PlayCount int    `json:"play_count" db:"play_count"`
}

type TrackPlayStatsResponse struct {
	Data  []TrackPlayStats `json:"data"`
	Total int              `json:"total"`
}

type ArtistPlayStatsResponse struct {
	Data []ArtistPlayStats `json:"data"`
}

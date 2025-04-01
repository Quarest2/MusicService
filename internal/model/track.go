package model

import "gorm.io/gorm"

type Track struct {
	gorm.Model
	Title      string `gorm:"not null"`
	Artist     string `gorm:"not null"`
	Album      string
	Genre      string
	Duration   int    // in seconds
	FilePath   string `gorm:"not null"` // path in MinIO
	UploadedBy uint   `gorm:"not null"` // user ID
}

type TrackUploadRequest struct {
	Title  string `form:"title" binding:"required"`
	Artist string `form:"artist" binding:"required"`
	Album  string `form:"album"`
	Genre  string `form:"genre"`
}

type TrackResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	Album     string `json:"album"`
	Genre     string `json:"genre"`
	Duration  int    `json:"duration"`
	CreatedAt string `json:"createdAt"`
}

type TrackSearchParams struct {
	Query  string `form:"q"`
	Artist string `form:"artist"`
	Album  string `form:"album"`
	Genre  string `form:"genre"`
}

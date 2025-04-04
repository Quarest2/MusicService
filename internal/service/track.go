package service

import (
	"MusicService/internal/model"
	"MusicService/internal/repository"
	"MusicService/internal/storage"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TrackService interface {
	UploadTrack(audioFile *multipart.FileHeader, imageFile *multipart.FileHeader, req *model.TrackUploadRequest, userID uint) (*model.TrackResponse, error)
	GetTrackByID(id uint) (*model.TrackResponse, error)
	GetAllTracks() ([]model.TrackResponse, error)
	StreamTrack(id uint) (io.ReadCloser, string, error)
	DeleteTrack(id uint) error
	SearchTracks(params model.TrackSearchParams) ([]model.TrackResponse, error)
	GetTrackImage(id uint) (io.ReadCloser, string, error)
	GetUserTracks(userId uint) ([]model.TrackResponse, error)
}

type trackService struct {
	trackRepo   repository.TrackRepository
	minioClient storage.MinIOClient
	bucketName  string
}

func NewTrackService(trackRepo repository.TrackRepository, minioClient storage.MinIOClient, bucketName string) TrackService {
	return &trackService{
		trackRepo:   trackRepo,
		minioClient: minioClient,
		bucketName:  bucketName,
	}
}

func (s *trackService) UploadTrack(audioFile *multipart.FileHeader, imageFile *multipart.FileHeader, req *model.TrackUploadRequest, userID uint) (*model.TrackResponse, error) {
	src, err := audioFile.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	ext := filepath.Ext(audioFile.Filename)
	audioFilename := uuid.New().String() + ext

	_, err = s.minioClient.PutObject(s.bucketName, audioFilename, src, audioFile.Size)
	if err != nil {
		return nil, err
	}

	src, err = imageFile.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	ext = filepath.Ext(imageFile.Filename)
	imageFilename := uuid.New().String() + ext

	_, err = s.minioClient.PutObject(s.bucketName, imageFilename, src, imageFile.Size)
	if err != nil {
		return nil, err
	}

	// TODO: In a real application, you'd want to extract duration and other metadata from the audio file
	track := &model.Track{
		Title:      req.Title,
		Artist:     req.Artist,
		Album:      req.Album,
		Genre:      req.Genre,
		Duration:   180, // TODO Placeholder - should extract from audio file
		FilePath:   audioFilename,
		ImagePath:  imageFilename,
		UploadedBy: userID,
	}

	if err := s.trackRepo.Create(track); err != nil {
		_ = s.minioClient.RemoveObject(s.bucketName, audioFilename)
		_ = s.minioClient.RemoveObject(s.bucketName, imageFilename)
		return nil, err
	}

	return &model.TrackResponse{
		ID:         track.ID,
		Title:      track.Title,
		Artist:     track.Artist,
		Album:      track.Album,
		Genre:      track.Genre,
		Duration:   track.Duration,
		CreatedAt:  track.CreatedAt.Format(time.RFC3339),
		UploadedBy: track.UploadedBy,
	}, nil
}

func (s *trackService) GetTrackByID(id uint) (*model.TrackResponse, error) {
	track, err := s.trackRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	var imageURL string
	if track.ImagePath != "" {
		imageURL = fmt.Sprintf("/api/tracks/%d/image", track.ID)
	}

	return &model.TrackResponse{
		ID:         track.ID,
		Title:      track.Title,
		Artist:     track.Artist,
		Album:      track.Album,
		Genre:      track.Genre,
		Duration:   track.Duration,
		ImageURL:   imageURL,
		CreatedAt:  track.CreatedAt.Format(time.RFC3339),
		UploadedBy: track.UploadedBy,
	}, nil
}

func (s *trackService) GetAllTracks() ([]model.TrackResponse, error) {
	tracks, err := s.trackRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var imageURL string

	var response []model.TrackResponse
	for _, track := range tracks {
		if track.ImagePath != "" {
			imageURL = fmt.Sprintf("/api/tracks/%d/image", track.ID)
		}

		response = append(response, model.TrackResponse{
			ID:         track.ID,
			Title:      track.Title,
			Artist:     track.Artist,
			Album:      track.Album,
			Genre:      track.Genre,
			Duration:   track.Duration,
			ImageURL:   imageURL,
			CreatedAt:  track.CreatedAt.Format(time.RFC3339),
			UploadedBy: track.UploadedBy,
		})
	}

	return response, nil
}

func (s *trackService) StreamTrack(id uint) (io.ReadCloser, string, error) {
	track, err := s.trackRepo.GetByID(id)
	if err != nil {
		return nil, "", err
	}

	object, err := s.minioClient.GetObject(s.bucketName, track.FilePath)
	if err != nil {
		return nil, "", err
	}

	contentType := "audio/mpeg"

	return object, contentType, nil
}

func (s *trackService) DeleteTrack(id uint) error {
	track, err := s.trackRepo.GetByID(id)
	if err != nil {
		return err
	}

	if err := s.minioClient.RemoveObject(s.bucketName, track.FilePath); err != nil {
		return err
	}

	return s.trackRepo.Delete(id)
}

func (s *trackService) SearchTracks(params model.TrackSearchParams) ([]model.TrackResponse, error) {
	tracks, err := s.trackRepo.Search(params)
	if err != nil {
		return nil, err
	}

	var response []model.TrackResponse
	for _, track := range tracks {
		response = append(response, model.TrackResponse{
			ID:         track.ID,
			Title:      track.Title,
			Artist:     track.Artist,
			Album:      track.Album,
			Genre:      track.Genre,
			Duration:   track.Duration,
			CreatedAt:  track.CreatedAt.Format(time.RFC3339),
			UploadedBy: track.UploadedBy,
		})
	}

	return response, nil
}

func (s *trackService) GetTrackImage(trackID uint) (io.ReadCloser, string, error) {
	track, err := s.trackRepo.GetByID(trackID)
	if err != nil {
		log.Printf("Error fetching track ID %d: %v", trackID, err)
		return nil, "", fmt.Errorf("failed to retrieve track: %w", err)
	}

	if track.ImagePath == "" {
		log.Printf("Track ID %d has no associated image", trackID)
		return nil, "", fmt.Errorf("track has no associated image")
	}

	obj, err := s.minioClient.GetObject(s.bucketName, track.ImagePath)
	if err != nil {
		log.Printf("MinIO error for path '%s': %v", track.ImagePath, err)
		return nil, "", fmt.Errorf("failed to retrieve image from storage: %w", err)
	}

	stat, err := obj.Stat()
	if err != nil {
		obj.Close()
		log.Printf("Invalid image object at '%s': %v", track.ImagePath, err)
		return nil, "", fmt.Errorf("image object verification failed: %w", err)
	}

	if stat.Size == 0 {
		obj.Close()
		log.Printf("Empty image file at '%s'", track.ImagePath)
		return nil, "", fmt.Errorf("image file is empty")
	}

	contentType := "application/octet-stream"
	switch {
	case strings.HasSuffix(track.ImagePath, ".jpg"),
		strings.HasSuffix(track.ImagePath, ".jpeg"):
		contentType = "image/jpeg"
	case strings.HasSuffix(track.ImagePath, ".png"):
		contentType = "image/png"
	case strings.HasSuffix(track.ImagePath, ".gif"):
		contentType = "image/gif"
	}

	log.Printf("Successfully retrieved image for track ID %d (%s)", trackID, track.ImagePath)
	return obj, contentType, nil
}

func (s *trackService) GetUserTracks(userId uint) ([]model.TrackResponse, error) {
	tracks, err := s.trackRepo.GetUserTracks(userId)
	if err != nil {
		return nil, err
	}

	var imageURL string

	var response []model.TrackResponse
	for _, track := range tracks {
		if track.ImagePath != "" {
			imageURL = fmt.Sprintf("/api/tracks/%d/image", track.ID)
		}

		response = append(response, model.TrackResponse{
			ID:         track.ID,
			Title:      track.Title,
			Artist:     track.Artist,
			Album:      track.Album,
			Genre:      track.Genre,
			Duration:   track.Duration,
			ImageURL:   imageURL,
			CreatedAt:  track.CreatedAt.Format(time.RFC3339),
			UploadedBy: track.UploadedBy,
		})
	}

	return response, nil
}

package service

import (
	"MusicService/internal/model"
	"MusicService/internal/repository"
	"MusicService/internal/storage"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type TrackService interface {
	UploadTrack(file *multipart.FileHeader, req *model.TrackUploadRequest, userID uint) (*model.TrackResponse, error)
	GetTrackByID(id uint) (*model.TrackResponse, error)
	GetAllTracks() ([]model.TrackResponse, error)
	StreamTrack(id uint) (io.ReadCloser, string, error)
	DeleteTrack(id uint) error
	SearchTracks(params model.TrackSearchParams) ([]model.TrackResponse, error)
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

func (s *trackService) UploadTrack(file *multipart.FileHeader, req *model.TrackUploadRequest, userID uint) (*model.TrackResponse, error) {
	// Open the file
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	newFilename := uuid.New().String() + ext

	// Upload to MinIO
	_, err = s.minioClient.PutObject(s.bucketName, newFilename, src, file.Size)
	if err != nil {
		return nil, err
	}

	// Create track record in database
	// Note: In a real application, you'd want to extract duration and other metadata from the audio file
	track := &model.Track{
		Title:      req.Title,
		Artist:     req.Artist,
		Album:      req.Album,
		Genre:      req.Genre,
		Duration:   180, // Placeholder - should extract from audio file
		FilePath:   newFilename,
		UploadedBy: userID,
	}

	if err := s.trackRepo.Create(track); err != nil {
		// Clean up MinIO object if database operation fails
		_ = s.minioClient.RemoveObject(s.bucketName, newFilename)
		return nil, err
	}

	return &model.TrackResponse{
		ID:        track.ID,
		Title:     track.Title,
		Artist:    track.Artist,
		Album:     track.Album,
		Genre:     track.Genre,
		Duration:  track.Duration,
		CreatedAt: track.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *trackService) GetTrackByID(id uint) (*model.TrackResponse, error) {
	track, err := s.trackRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &model.TrackResponse{
		ID:        track.ID,
		Title:     track.Title,
		Artist:    track.Artist,
		Album:     track.Album,
		Genre:     track.Genre,
		Duration:  track.Duration,
		CreatedAt: track.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *trackService) GetAllTracks() ([]model.TrackResponse, error) {
	tracks, err := s.trackRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var response []model.TrackResponse
	for _, track := range tracks {
		response = append(response, model.TrackResponse{
			ID:        track.ID,
			Title:     track.Title,
			Artist:    track.Artist,
			Album:     track.Album,
			Genre:     track.Genre,
			Duration:  track.Duration,
			CreatedAt: track.CreatedAt.Format(time.RFC3339),
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

	// Get content type (simplified - in real app you'd detect from file)
	contentType := "audio/mpeg"

	return object, contentType, nil
}

func (s *trackService) DeleteTrack(id uint) error {
	track, err := s.trackRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Delete from MinIO first
	if err := s.minioClient.RemoveObject(s.bucketName, track.FilePath); err != nil {
		return err
	}

	// Then delete from database
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
			ID:        track.ID,
			Title:     track.Title,
			Artist:    track.Artist,
			Album:     track.Album,
			Genre:     track.Genre,
			Duration:  track.Duration,
			CreatedAt: track.CreatedAt.Format(time.RFC3339),
		})
	}

	return response, nil
}

package storage

import (
	"MusicService/internal/config"
	"context"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOClient интерфейс для работы с MinIO
type MinIOClient interface {
	CreateBucket(bucketName string) error
	PutObject(bucketName, objectName string, reader io.Reader, objectSize int64) (minio.UploadInfo, error)
	GetObject(bucketName, objectName string) (*minio.Object, error)
	RemoveObject(bucketName, objectName string) error
	PresignedGetObject(bucketName, objectName string, expiry time.Duration) (*url.URL, error)
}

type minioClient struct {
	client *minio.Client
}

// NewMinioClient создает новый клиент MinIO
func NewMinioClient(cfg *config.Config) (MinIOClient, error) {
	// Инициализация клиента MinIO
	client, err := minio.New(cfg.MinIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIO.AccessKey, cfg.MinIO.SecretKey, ""),
		Secure: cfg.MinIO.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	return &minioClient{client: client}, nil
}

// CreateBucket создает bucket в MinIO если он не существует
func (m *minioClient) CreateBucket(bucketName string) error {
	ctx := context.Background()
	exists, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if !exists {
		err = m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}

		// Устанавливаем политику доступа (по умолчанию private)
		policy := `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::` + bucketName + `/*"]
				}
			]
		}`
		err = m.client.SetBucketPolicy(ctx, bucketName, policy)
		if err != nil {
			return err
		}
	}

	return nil
}

// PutObject загружает объект в MinIO
func (m *minioClient) PutObject(bucketName, objectName string, reader io.Reader, objectSize int64) (minio.UploadInfo, error) {
	uploadInfo, err := m.client.PutObject(
		context.Background(),
		bucketName,
		objectName,
		reader,
		objectSize,
		minio.PutObjectOptions{
			ContentType: "application/octet-stream",
		},
	)
	return uploadInfo, err
}

// GetObject получает объект из MinIO
func (m *minioClient) GetObject(bucketName, objectName string) (*minio.Object, error) {
	object, err := m.client.GetObject(
		context.Background(),
		bucketName,
		objectName,
		minio.GetObjectOptions{},
	)
	return object, err
}

// RemoveObject удаляет объект из MinIO
func (m *minioClient) RemoveObject(bucketName, objectName string) error {
	err := m.client.RemoveObject(
		context.Background(),
		bucketName,
		objectName,
		minio.RemoveObjectOptions{},
	)
	return err
}

// PresignedGetObject генерирует временную ссылку для доступа к объекту
func (m *minioClient) PresignedGetObject(bucketName, objectName string, expiry time.Duration) (*url.URL, error) {
	url, err := m.client.PresignedGetObject(
		context.Background(),
		bucketName,
		objectName,
		expiry,
		nil,
	)
	return url, err
}

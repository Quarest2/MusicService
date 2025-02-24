package minioStorage

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

type Client interface {
	InitMinio() error
	CreateOne(file FileDataType) (string, error)
	CreateMany(map[string]FileDataType) ([]string, error)
	GetOne(objectID string) (string, error)
	GetMany(objectIDs []string) ([]string, error)
	DeleteOne(objectID string) error
	DeleteMany(objectIDs []string) error
}

type minioClient struct {
	mc *minio.Client
}

func NewMinioClient() Client {
	return &minioClient{}
}

func (m *minioClient) InitMinio() error {
	var err error

	ctx := context.Background()

	var client *minio.Client
	if client, err = minio.New(AppConfig.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(AppConfig.MinioRootUser, AppConfig.MinioRootPassword, ""),
		Secure: AppConfig.MinioUseSSL,
	}); err != nil {
		log.Printf("Failed to create minio client: %v", err)
		return err
	}

	m.mc = client

	var exists bool
	if exists, err = m.mc.BucketExists(ctx, AppConfig.BucketName); err != nil {
		log.Printf("Failed to check if bucket exists: %v", err)
		return err
	}

	if !exists {
		if err = m.mc.MakeBucket(ctx, AppConfig.BucketName, minio.MakeBucketOptions{}); err != nil {
			log.Printf("Failed to create bucket: %v", err)
			return err
		}
	}

	return nil
}

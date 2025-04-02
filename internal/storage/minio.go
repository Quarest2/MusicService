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

func NewMinioClient(cfg *config.Config) (MinIOClient, error) {
	client, err := minio.New(cfg.MinIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIO.AccessKey, cfg.MinIO.SecretKey, ""),
		Secure: cfg.MinIO.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	return &minioClient{client: client}, nil
}

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

func (m *minioClient) GetObject(bucketName, objectName string) (*minio.Object, error) {
	object, err := m.client.GetObject(
		context.Background(),
		bucketName,
		objectName,
		minio.GetObjectOptions{},
	)
	return object, err
}

func (m *minioClient) RemoveObject(bucketName, objectName string) error {
	err := m.client.RemoveObject(
		context.Background(),
		bucketName,
		objectName,
		minio.RemoveObjectOptions{},
	)
	return err
}

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

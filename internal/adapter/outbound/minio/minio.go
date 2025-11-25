package minio

import (
	"bytes"
	"clean-architecture/internal/port/outbound"
	"context"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

type MinioStorage struct {
	Client     *minio.Client
	BucketName string
}

func NewMinioStorage(client *minio.Client, bucket string) outbound.MinioInterface {
	return &MinioStorage{
		Client:     client,
		BucketName: bucket,
	}
}

func (m *MinioStorage) UploadFile(path string, file *bytes.Buffer) (string, error) {
	_, err := m.Client.PutObject(
		context.Background(),
		m.BucketName,
		path,
		bytes.NewReader(file.Bytes()),
		int64(file.Len()),
		minio.PutObjectOptions{
			ContentType: "image/jpeg", // bisa dinamis
		},
	)
	if err != nil {
		return "", err
	}

	return path, nil
}

func (m *MinioStorage) GetPresignedURL(path string, expiry time.Duration) (string, error) {
	reqParams := make(url.Values)

	presignedURL, err := m.Client.PresignedGetObject(
		context.Background(),
		m.BucketName,
		path,
		expiry,
		reqParams,
	)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

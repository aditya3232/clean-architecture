package minio

import (
	"bytes"
	"clean-architecture/internal/port/outbound"
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
)

type MinioStorage struct {
	Client     *minio.Client
	BucketName string
	BaseURL    string // domain minio yang bisa diakses client
}

func NewMinioStorage(client *minio.Client, bucket, baseUrl string) outbound.MinioInterface {
	return &MinioStorage{
		Client:     client,
		BucketName: bucket,
		BaseURL:    baseUrl,
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

	// URL publik / signed URL sesuai kebutuhan
	url := fmt.Sprintf("%s/%s/%s", m.BaseURL, m.BucketName, path)
	return url, nil
}

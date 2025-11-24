package config

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
)

func (cfg Config) InitMinio() (*minio.Client, error) {
	minioClient, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.AccessKey, cfg.Minio.SecretKey, ""),
		Secure: cfg.Minio.UseSSL,
	})

	if err != nil {
		log.Error().Err(err).Msg("[ConnectionMinio-1] Failed to connect to storage " + cfg.Minio.Endpoint)
		return nil, err
	}

	_, err = minioClient.ListBuckets(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("[ConnectionMinio-2] Failed to get list bucket")
		return nil, err
	}

	return minioClient, nil

}

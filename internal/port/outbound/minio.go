package outbound

import (
	"bytes"
	"time"
)

type MinioInterface interface {
	UploadFile(path string, file *bytes.Buffer) (string, error)
	GetPresignedURL(path string, expiry time.Duration) (string, error)
}

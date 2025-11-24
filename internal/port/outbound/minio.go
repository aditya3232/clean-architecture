package outbound

import "bytes"

type MinioInterface interface {
	UploadFile(path string, file *bytes.Buffer) (string, error)
}

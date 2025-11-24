package echo

import (
	"bytes"
	"clean-architecture/internal/adapter/inbound/echo/response"
	"clean-architecture/internal/port/inbound"
	"clean-architecture/internal/port/outbound"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type uploadImageHandler struct {
	storageHandler outbound.MinioInterface
}

func NewUploadImageHandler(storageHandler outbound.MinioInterface) inbound.UploadImageInterface {
	return &uploadImageHandler{storageHandler: storageHandler}
}

func (u *uploadImageHandler) UploadImage(c echo.Context) error {
	var resp = response.DefaultResponse{}

	file, err := c.FormFile("photo")
	if err != nil {
		return response.RespondWithError(c, http.StatusUnprocessableEntity, "[UploadImage-1] UploadImage", err)
	}

	src, err := file.Open()
	if err != nil {
		return response.RespondWithError(c, http.StatusBadRequest, "[UploadImage-2] UploadImage", err)
	}

	defer func() {
		if cerr := src.Close(); cerr != nil {
			log.Errorf("[UploadImage-defer] failed to close src: %v", cerr)
		}
	}()

	fileBuffer := new(bytes.Buffer)
	_, err = io.Copy(fileBuffer, src)
	if err != nil {
		return response.RespondWithError(c, http.StatusInternalServerError, "[UploadImage-4] UploadImage", err)
	}

	newFileName := fmt.Sprintf("%s_%d%s",
		uuid.New().String(),
		time.Now().Unix(),
		getExtension(file.Filename),
	)

	uploadPath := fmt.Sprintf("public/uploads/%s", newFileName)

	url, err := u.storageHandler.UploadFile(uploadPath, fileBuffer)
	if err != nil {
		return response.RespondWithError(c, http.StatusInternalServerError, "[UploadImage-5] UploadImage", err)
	}

	resp.Message = "Success"
	resp.Data = map[string]string{"image_url": url}

	return c.JSON(http.StatusOK, resp)
}

func getExtension(fileName string) string {
	ext := "." + fileName[len(fileName)-3:] // Ambil 3 karakter terakhir untuk ekstensi
	if len(fileName) > 4 && fileName[len(fileName)-4] == '.' {
		ext = "." + fileName[len(fileName)-4:]
	}
	return ext
}

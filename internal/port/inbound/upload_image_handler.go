package inbound

import "github.com/labstack/echo/v4"

type UploadImageInterface interface {
	UploadImage(c echo.Context) error
}

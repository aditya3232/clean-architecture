package inbound

import "github.com/labstack/echo/v4"

type RoleHandlerInterface interface {
	GetAll(c echo.Context) error
	GetByID(c echo.Context) error
	Create(c echo.Context) error
	Delete(c echo.Context) error
	Update(c echo.Context) error
}

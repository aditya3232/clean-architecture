package inbound

import "github.com/labstack/echo/v4"

type UserHandlerInterface interface {
	SignIn(c echo.Context) error
	CreateUserAccount(c echo.Context) error
	ForgotPassword(c echo.Context) error
	VerifyAccount(c echo.Context) error
	UpdatePassword(c echo.Context) error
	GetProfileUser(c echo.Context) error
	UpdateDataUser(c echo.Context) error

	// Modul Customers Admin
	GetCustomerAll(c echo.Context) error
	GetCustomerByID(c echo.Context) error
	CreateCustomer(c echo.Context) error
	UpdateCustomer(c echo.Context) error
	DeleteCustomer(c echo.Context) error
}

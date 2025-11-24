package echo

import (
	"clean-architecture/internal/port/inbound"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRoutes(
	e *echo.Echo,
	mid inbound.MiddlewareAdapterInterface,
	pingHandler inbound.PingHandlerInterface,
	userHandler inbound.UserHandlerInterface,
	roleHandler inbound.RoleHandlerInterface,
	uploadImageHandler inbound.UploadImageInterface,
) {
	e.Use(middleware.Recover())

	e.GET("/ping", pingHandler.Ping)

	e.POST("/signin", userHandler.SignIn)
	e.POST("/signup", userHandler.CreateUserAccount)
	e.POST("/forgot-password", userHandler.ForgotPassword)
	e.GET("/verify-account", userHandler.VerifyAccount)
	e.PUT("/update-password", userHandler.UpdatePassword)

	adminGroup := e.Group("/admin", mid.CheckToken())
	adminGroup.GET("/customers", userHandler.GetCustomerAll)
	adminGroup.POST("/customers", userHandler.CreateCustomer)
	adminGroup.PUT("/customers/:id", userHandler.UpdateCustomer)
	adminGroup.GET("/customers/:id", userHandler.GetCustomerByID)
	adminGroup.DELETE("/customers/:id", userHandler.DeleteCustomer)

	adminGroup.GET("/roles", roleHandler.GetAll)
	adminGroup.POST("/roles", roleHandler.Create)
	adminGroup.PUT("/roles/:id", roleHandler.Update)
	adminGroup.DELETE("/roles/:id", roleHandler.Delete)
	adminGroup.GET("/roles/:id", roleHandler.GetByID)

	authGroup := e.Group("/auth", mid.CheckToken())
	authGroup.GET("/profile", userHandler.GetProfileUser)
	authGroup.PUT("/profile", userHandler.UpdateDataUser)
	authGroup.POST("/profile/image-upload", uploadImageHandler.UploadImage)
}

package handler

import (
	"RPJ-Overseas-Exim/yourpharma-admin/handler/admin-handler"
	"RPJ-Overseas-Exim/yourpharma-admin/handler/auth-handler"

	"github.com/labstack/echo/v4"
)

func SetupAuthRoutes(e *echo.Echo) {
	e.GET("/", authHandler.LoginHandler)
	e.POST("/", authHandler.LoginHandler)
	e.GET("/register", authHandler.RegisterHandler)
}

func SetupAdminRoutes(e *echo.Echo) {
	e.GET("/home", adminHandler.Home)
	e.GET("/customers", adminHandler.Customers)
	e.GET("/orders", adminHandler.Orders)
}

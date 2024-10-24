package handlers

import (
	"RPJ-Overseas-Exim/yourpharma-admin/handler/admin-handler"
	"RPJ-Overseas-Exim/yourpharma-admin/handler/auth-handler"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type handler struct{
    DB *gorm.DB
}

func New(db *gorm.DB) handler {
    return handler{db}
}

func (h *handler) SetupAuthRoutes(e *echo.Echo) {
    as := authHandler.NewAuthService(h.DB)
	e.GET("/", as.LoginHandler)
	e.POST("/", as.LoginHandler)

	e.GET("/register", as.RegisterHandler)
}

func (h *handler) SetupHomeRoutes(e *echo.Echo) {
    hs := adminHandler.NewHomeService(h.DB)
	e.GET("/home", hs.Home)
	e.POST("/home", hs.Home)
}

func (h *handler) SetupCustomerRoutes(e *echo.Echo) {
    cs := adminHandler.NewCustomerService(h.DB)
	e.GET("/customers", cs.Customers)
	e.POST("/customers", cs.CreateCustomer)
}

func (h *handler) SetupOrderRoutes(e *echo.Echo) {
    ords := adminHandler.NewOrderService(h.DB)
	e.GET("/orders", ords.Orders)
}

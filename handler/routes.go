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

func (h *handler) SetupCustomerRoutes(e *echo.Echo) {
    cs := adminHandler.NewCustomerService(h.DB)
	e.GET("/customers", cs.Customers)
	e.POST("/customers", cs.CreateCustomer)
	e.PUT("/customers", cs.UpdateCustomer)
    e.DELETE("/customers/:id", cs.DeleteCustomer)
}

func (h *handler) SetupProductRoutes(e *echo.Echo){
    ps := adminHandler.NewProductService(h.DB)
    e.GET("/products", ps.Products)
    e.POST("/products", ps.CreateProduct)
    e.PUT("/products/:id", ps.UpdateProduct)
    e.DELETE("/price/:id", ps.DeletePrice)
    e.DELETE("/products/:id", ps.DeleteProduct)
}

func (h *handler) SetupOrderRoutes(e *echo.Echo) {
    ords := adminHandler.NewOrderService(h.DB)
	e.GET("/orders", ords.Orders)
	e.POST("/orders", ords.CreateOrder)
    e.PUT("/orders/:id", ords.UpdateOrder)
    e.DELETE("/orders/:id", ords.DeleteOrder)
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



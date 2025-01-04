package handlers

import (
	"RPJ-Overseas-Exim/yourpharma-admin/handler/admin-handler"
	"RPJ-Overseas-Exim/yourpharma-admin/handler/auth-handler"
	authMiddleware "RPJ-Overseas-Exim/yourpharma-admin/pkg/middleware"

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

    customerRoute := e.Group("/customers", authMiddleware.AuthMiddleware)
	customerRoute.GET("", cs.Customers)
	customerRoute.POST("", cs.CreateCustomer)
    customerRoute.PUT("/:id", cs.UpdateCustomer)
    customerRoute.DELETE("/:id", cs.DeleteCustomer)
    customerRoute.POST("/import", cs.ImportCustomers)
}

func (h *handler) SetupProductRoutes(e *echo.Echo){
    ps := adminHandler.NewProductService(h.DB)

    productRoute := e.Group("/products", authMiddleware.AuthMiddleware)
    productRoute.GET("", ps.Products)
    productRoute.POST("", ps.CreateProduct)
    productRoute.PUT("/:id", ps.UpdateProduct)
    productRoute.DELETE("/:id", ps.DeleteProduct)

    priceRoute := e.Group("/price", authMiddleware.AuthMiddleware)
    priceRoute.DELETE("/:id", ps.DeletePrice)
    priceRoute.POST("", ps.CreatePrice)
}

func (h *handler) SetupOrderRoutes(e *echo.Echo) {
    ords := adminHandler.NewOrderService(h.DB)

    orderRoute := e.Group("/orders", authMiddleware.AuthMiddleware)
	orderRoute.GET("", ords.Orders)
	orderRoute.POST("", ords.CreateOrder)
    orderRoute.PUT("/:id", ords.UpdateOrder)
    orderRoute.DELETE("/:id", ords.DeleteOrder)
    orderRoute.POST("/import", ords.ImportOrders)
}

func (h *handler) SetupAuthRoutes(e *echo.Echo) {
    as := authHandler.NewAuthService(h.DB)
	e.GET("/", as.LoginHandler, authMiddleware.LoginMiddleware)
	e.POST("/", as.LoginHandler, authMiddleware.LoginMiddleware)
}

func (h *handler) SetupHomeRoutes(e *echo.Echo) {
    hs := adminHandler.NewHomeService(h.DB)

    homeRoute := e.Group("/home", authMiddleware.AuthMiddleware)
	homeRoute.GET("", hs.Home)
	homeRoute.POST("", hs.Home)
}



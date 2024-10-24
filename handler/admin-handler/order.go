package adminHandler

import (
	authHandler "RPJ-Overseas-Exim/yourpharma-admin/handler/auth-handler"
	adminView "RPJ-Overseas-Exim/yourpharma-admin/templ/admin-views"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type orderService struct {
    DB *gorm.DB
}

func NewOrderService(db *gorm.DB) *orderService{
    return &orderService{DB: db}
}

func (ords *orderService) Orders(c echo.Context) error {
	ordersView := adminView.Orders()
	var msgs []string

	return authHandler.RenderView(c, adminView.AdminIndex(
		"Orders",
		true,
		msgs,
		msgs,
		ordersView,
	))
}

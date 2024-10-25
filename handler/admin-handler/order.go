package adminHandler

import (
	authHandler "RPJ-Overseas-Exim/yourpharma-admin/handler/auth-handler"
	"RPJ-Overseas-Exim/yourpharma-admin/pkg/types"
	adminView "RPJ-Overseas-Exim/yourpharma-admin/templ/admin-views"
	"log"

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
    var customerData []types.Order
    err := ords.DB.Find(&customerData)

    if err!=nil{
        log.Printf("Customers not present: %v", err)
    }

	ordersView := adminView.Orders(customerData)
	var msgs []string

	return authHandler.RenderView(c, adminView.AdminIndex(
		"Orders",
		true,
		msgs,
		msgs,
		ordersView,
	))
}

package adminHandler

import (
	authHandler "RPJ-Overseas-Exim/yourpharma-admin/handler/auth-handler"
	"RPJ-Overseas-Exim/yourpharma-admin/pkg/types"
	"RPJ-Overseas-Exim/yourpharma-admin/pkg/utils"
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

// database functions ===================================================
func (ords *orderService) GetOrders() ([]types.Order, error) {
    var ordersData []types.Order

    result := ords.DB.Find(&ordersData)
    if result.Error != nil {
        return ordersData, result.Error
    }

    return ordersData, nil
}

func (ords *orderService) AddOrderDetails() error {
    return nil
}

func (ords *orderService) UpdataeOrderDetails() error {
    return nil
}

func (ords *orderService) DeleteOrderDetails() error {
    return nil
}


// routes functions ===================================================
func (ords *orderService) Orders(c echo.Context) error {
    var err error
    ordersData, err := ords.GetOrders()
    utils.ErrorHandler(err, "Failed to get the order data")

	ordersView := adminView.Orders(ordersData)
	var msgs []string

	return authHandler.RenderView(c, adminView.AdminIndex(
		"Orders",
		true,
		msgs,
		msgs,
		ordersView,
	))
}

func (ords *orderService) CreateOrder(c echo.Context) error {
    var err error



    ordersData, err := ords.GetOrders()
    utils.ErrorHandler(err, "Failed to get the order data")
    orderView := adminView.Orders(ordersData)
    return authHandler.RenderView(c, orderView)
}

func (ords *orderService) UpdateOrder(c echo.Context) error {
    var err error
    ordersData, err := ords.GetOrders()
    utils.ErrorHandler(err, "Failed to get the order data")
    orderView := adminView.Orders(ordersData)
    return authHandler.RenderView(c, orderView)
}

func (ords *orderService) DeleteOrder(c echo.Context) error {
    var err error
    ordersData, err := ords.GetOrders()
    utils.ErrorHandler(err, "Failed to get the order data")
    orderView := adminView.Orders(ordersData)
    return authHandler.RenderView(c, orderView)
}

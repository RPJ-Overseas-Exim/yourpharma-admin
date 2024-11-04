package adminHandler

import (
	"RPJ-Overseas-Exim/yourpharma-admin/db/models"
	"RPJ-Overseas-Exim/yourpharma-admin/handler/auth-handler"
	"RPJ-Overseas-Exim/yourpharma-admin/pkg/types"
	"RPJ-Overseas-Exim/yourpharma-admin/templ/admin-views"
	"log"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type homeService struct{
    DB *gorm.DB
}

func NewHomeService(db *gorm.DB) *homeService {
    return &homeService{DB: db}
}

// data base functions for home ================================================
func (hs *homeService) GetOrders(status string) ([]types.Order, error) {
    var ordersData []types.Order
    var result *gorm.DB
    selectStatement := `
                orders.id as Id,
                customers.name as Name,
                customers.email as Email,
                customers.number as Number,
                customers.address as Address,
                products.name as Product,
                orders.quantity as Quantity,
                orders.status as Status,
                orders.origin as Origin,
                orders.amount as Price,
                orders.created_at as CreatedAt,
                orders.updated_at as UpdatedAt
            `
    joinStatement1 := `
                inner join customers
                on customers.id = orders.customer_id
            `
    joinStatement2 := `
                inner join products
                on products.id = orders.product_id
            `

    if len(status) != 0 && status != "all" {
        result = hs.DB.Model(&models.Order{}).Select(selectStatement).Joins(joinStatement1).Joins(joinStatement2).Where("orders.status like ?", status).Scan(&ordersData)
    }else{
        result = hs.DB.Model(&models.Order{}).Select(selectStatement).Joins(joinStatement1).Joins(joinStatement2).Scan(&ordersData)
    }

    log.Printf("Order result: %v", ordersData)

    if result.Error != nil {
        return ordersData, result.Error
    }

    return ordersData, nil
}

// routes functions for home ===================================================
func (hs *homeService) Home(c echo.Context) error {
    ordersData, err := hs.GetOrders("active")
    if err!=nil {
        log.Printf("Failed to get the orders data")
    }
	homeView := adminView.Home(29000, 405, 24, 300, 320, ordersData)

	return authHandler.RenderView(c, adminView.AdminIndex(
		"Home",
		true,
		homeView,
	))
}


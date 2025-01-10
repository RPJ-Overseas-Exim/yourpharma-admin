package adminHandler

import (
	"RPJ-Overseas-Exim/yourpharma-admin/db/models"
	"RPJ-Overseas-Exim/yourpharma-admin/handler/auth-handler"
	"RPJ-Overseas-Exim/yourpharma-admin/pkg/types"
	"RPJ-Overseas-Exim/yourpharma-admin/pkg/utils"
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
        result = hs.DB.Model(&models.Order{}).Select(selectStatement).Joins(joinStatement1).Joins(joinStatement2).Where("orders.status ilike ?", status).Limit(10).Scan(&ordersData)
    }else{
        result = hs.DB.Model(&models.Order{}).Select(selectStatement).Joins(joinStatement1).Joins(joinStatement2).Scan(&ordersData)
    }

    log.Printf("Order result: %v", ordersData)

    if result.Error != nil {
        return ordersData, result.Error
    }

    return ordersData, nil
}

func (hs *homeService) GetTotalSales() int {
    var totalSales int

    result := hs.DB.Model(&models.Order{}).Select("sum(amount) as Total").Where("status ilike ?", "delivered").Find(&totalSales)
    if result.Error != nil {
        return 0
    }

    return totalSales
}

func (hs *homeService) GetTotalOrders() int {
    var totalOrders int

    result := hs.DB.Model(&models.Order{}).Select("count(amount) as Total").Find(&totalOrders)
    if result.Error != nil {
        return 0
    }

    return totalOrders
}

func (hs *homeService) GetTotalOrderInProcess() int {
    var totalOrderInProcess int

    result := hs.DB.Model(&models.Order{}).Select("count(amount) as total").Where("status ilike ? or status ilike ?", "paid", "shipped").Find(&totalOrderInProcess)
    if result.Error != nil {
        return 0
    }

    return totalOrderInProcess
}

func (hs *homeService) GetTotalOrderDelivered() int {
    var totalOrderDelivered int

    result := hs.DB.Model(&models.Order{}).Select("count(amount) as total").Where("status ilike ?", "delivered").Find(&totalOrderDelivered)
    if result.Error != nil {
        return 0
    }

    return totalOrderDelivered
}

func (hs *homeService) GetTotalCustomers() int {
    var totalCustomers int

    result := hs.DB.Model(&models.Customer{}).Select("count(email) as total").Find(&totalCustomers)
    if result.Error != nil {
        return 0
    }

    return totalCustomers
}

// routes functions for home ===================================================
func (hs *homeService) Home(c echo.Context) error {
    ordersData, err := hs.GetOrders("active")
    if err!=nil {
        log.Printf("Failed to get the orders data")
    }

    totalSales := hs.GetTotalSales()
    totalOrders := hs.GetTotalOrders()
    totalOrderInProcess := hs.GetTotalOrderInProcess()
    totalOrderDelivered := hs.GetTotalOrderDelivered()
    totalCustomers := hs.GetTotalCustomers()

	homeView := adminView.Home(totalSales, totalOrders, totalOrderInProcess, totalOrderDelivered, totalCustomers, ordersData)

    role := utils.GetRole(utils.GetAdmin(c))
	return authHandler.RenderView(c, adminView.AdminIndex(
		"Home",
		true,
		homeView,
        role,
	))
}


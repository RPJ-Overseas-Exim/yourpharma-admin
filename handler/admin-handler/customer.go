package adminHandler

import (
	authHandler "RPJ-Overseas-Exim/yourpharma-admin/handler/auth-handler"
	"RPJ-Overseas-Exim/yourpharma-admin/pkg/types"
	adminView "RPJ-Overseas-Exim/yourpharma-admin/templ/admin-views"
	"log"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Customer struct{
    Name string
    Email string
    Number string
    Address string
}

type customerService struct{
    DB *gorm.DB
}

func NewCustomerService(db *gorm.DB) *customerService {
    return &customerService{DB: db}
}

func (cs *customerService) Customers(c echo.Context) error {
    var customerData []types.Customer
    result := cs.DB.Find(&customerData)

    if result.Error != nil {
        log.Printf("Error in customers Data: %v", result.Error)
        return c.JSON(400, &echo.Map{"message":"Customers not present"})
    }

	customersView := adminView.Customers(customerData)
	var msgs []string

	return authHandler.RenderView(c, adminView.AdminIndex(
		"Customers",
		true,
		msgs,
		msgs,
		customersView,
	))
}

func (cs *customerService) CreateCustomer(c echo.Context) error {
    var customer types.Customer
    var result *gorm.DB

    customer = types.Customer{
        Name: c.FormValue("name"),
        Email: c.FormValue("email"),
        Number: c.FormValue("number"),
        Address: c.FormValue("address"),
    }

    result = cs.DB.Create(&customer)

    var customersData []types.Customer
    result = cs.DB.Find(&customersData) 

    if result.Error != nil {
        log.Printf("Customers data is not fetched: %v", result.Error)
    }

    customersView := adminView.Customers(customersData)
    return authHandler.RenderView(c, customersView)
}

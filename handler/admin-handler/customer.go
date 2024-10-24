package adminHandler

import (
	authHandler "RPJ-Overseas-Exim/yourpharma-admin/handler/auth-handler"
	handlerUtils "RPJ-Overseas-Exim/yourpharma-admin/handler/utils"
	adminView "RPJ-Overseas-Exim/yourpharma-admin/templ/admin-views"

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
	customersView := adminView.Customers()
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
    var err error
    var data map[string]interface{}

    data, err = handlerUtils.ResponseBody(c)

    if err != nil {
        return c.JSON(400, &echo.Map{"message":"Data is not provided"})
    }

    customer := Customer{
        Name: data["name"].(string),
        Email: data["email"].(string),
        Number: data["number"].(string),
        Address: data["address"].(string),
    }

    result := cs.DB.Create(&customer)
    if result.Error != nil {
        return c.JSON(401, &echo.Map{"message":"failed to create customer or customer is already present"})
    }

    return c.JSON(200, &echo.Map{"message": "Customer created"})
}

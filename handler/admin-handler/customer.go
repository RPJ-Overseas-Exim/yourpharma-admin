package adminHandler

import (
	"RPJ-Overseas-Exim/yourpharma-admin/db/models"
	authHandler "RPJ-Overseas-Exim/yourpharma-admin/handler/auth-handler"
	"RPJ-Overseas-Exim/yourpharma-admin/pkg/types"
	adminView "RPJ-Overseas-Exim/yourpharma-admin/templ/admin-views"
	"log"
	"strconv"

	"github.com/aidarkhanov/nanoid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type customerService struct{
    DB *gorm.DB
}

func NewCustomerService(db *gorm.DB) *customerService {
    return &customerService{DB: db}
}

// customers data methods
func (cs *customerService) GetCustomers() ([]types.Customer, error) {
    var customersData []types.Customer
    result := cs.DB.Find(&customersData, "deleted_at is NULL")

    if result.Error != nil {
        return customersData, result.Error
    }

    return customersData, nil
}

func (cs *customerService) AddCustomer(name, email string, number *int, address string) error {
    id := nanoid.New()
    customer := types.Customer{
        Id: id,
        Name: name,
        Email: email,
        Number: number,
        Address: address,
    }

    err := cs.DB.Create(&customer).Error

    if err != nil {
        return  err
    }
    return nil
}

func (cs *customerService) UpdateCustomerDetails(id, name, email string, number *int, address string) error {
    result := cs.DB.Model(models.Customer{}).Where("id = ?", id).Updates(map[string]interface{}{
        "name": name, 
        "email": email,
        "number": number,
        "address": address,
    })

    if result.Error != nil {
        return result.Error
    }
    return nil
}

func (cs *customerService) DeleteCustomerDetails(id string) error {
    result := cs.DB.Delete(&models.Customer{}, "id like ?", id)

    if result.Error != nil {
        log.Printf("Failed to delete the customer: %v", result.Error)
        return result.Error
    }

    return nil 
}


// routes methods
func (cs *customerService) Customers(c echo.Context) error {
    customersData, err := cs.GetCustomers()

    if err != nil {
        log.Printf("Error in customers Data: %v", err)
        return c.JSON(400, &echo.Map{"message":"Customers not present"})
    }

	customersView := adminView.Customers(customersData)
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
    num, err :=  strconv.Atoi(c.FormValue("number"))
    err = cs.AddCustomer(c.FormValue("name"), c.FormValue("email"), &num, c.FormValue("address"))

    customersData, err := cs.GetCustomers()
    if err != nil {
        log.Printf("Customers data is not fetched: %v", err)
    }

    customersView := adminView.Customers(customersData)
    return authHandler.RenderView(c, customersView)
}

func (cs *customerService) UpdateCustomer(c echo.Context) error {
    var err error
    num, err := strconv.Atoi(c.FormValue("number"))
    err = cs.UpdateCustomerDetails(
        c.FormValue("id"),
        c.FormValue("name"),
        c.FormValue("email"),
        &num,
        c.FormValue("address"),
    )

    if err != nil{
        return err
    }

    customersData, err := cs.GetCustomers()

    customerView := adminView.Customers(customersData)
    return authHandler.RenderView(c, customerView)
}

func (cs *customerService) DeleteCustomer(c echo.Context) error {
    var err error
    err = cs.DeleteCustomerDetails(c.Param("id"))

    if err != nil {
        return err
    }
    
    customersData, err := cs.GetCustomers()

    customerView := adminView.Customers(customersData)
    return authHandler.RenderView(c, customerView)
}

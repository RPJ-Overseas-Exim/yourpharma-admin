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
func (cs *customerService) GetCustomers(page int, limit int) ([]types.Customer, error) {
    var customersData []types.Customer
    result := cs.DB.Find(&customersData, "deleted_at is NULL").Offset(page*limit).Limit(limit)

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
    var err error
    page, err := strconv.Atoi(c.QueryParam("page"))
    if err!=nil {
        page = 0
    }
    limit, err := strconv.Atoi(c.QueryParam("limit"))
    if err != nil{
        limit = 10
    }
    customersData, err := cs.GetCustomers(page, limit)

    if err != nil {
        log.Printf("Error in customers Data: %v", err)
        return c.JSON(400, &echo.Map{"message":"Customers not present"})
    }

	customersView := adminView.Customers(customersData)

	return authHandler.RenderView(c, adminView.AdminIndex(
		"Customers",
		true,
		customersView,
	))
}

func (cs *customerService) CreateCustomer(c echo.Context) error {
    var err error
    num, err :=  strconv.Atoi(c.FormValue("number"))
    err = cs.AddCustomer(c.FormValue("name"), c.FormValue("email"), &num, c.FormValue("address"))

    page, err := strconv.Atoi(c.QueryParam("page"))
    if err!=nil {
        page = 0
    }
    limit, err := strconv.Atoi(c.QueryParam("limit"))
    if err != nil{
        limit = 10
    }

    customersData, err := cs.GetCustomers(page, limit)
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
        c.Param("id"),
        c.FormValue("name"),
        c.FormValue("email"),
        &num,
        c.FormValue("address"),
    )

    if err != nil{
        return err
    }

    page, err := strconv.Atoi(c.QueryParam("page"))
    if err!=nil {
        page = 0
    }
    limit, err := strconv.Atoi(c.QueryParam("limit"))
    if err != nil{
        limit = 10
    }

    customersData, err := cs.GetCustomers(page, limit)
    customerView := adminView.Customers(customersData)
    return authHandler.RenderView(c, customerView)
}

func (cs *customerService) DeleteCustomer(c echo.Context) error {
    var err error
    err = cs.DeleteCustomerDetails(c.Param("id"))

    if err != nil {
        return err
    }

    page, err := strconv.Atoi(c.QueryParam("page"))
    if err!=nil {
        page = 0
    }
    limit, err := strconv.Atoi(c.QueryParam("limit"))
    if err != nil{
        limit = 10
    }

    customersData, err := cs.GetCustomers(page, limit)
    customerView := adminView.Customers(customersData)
    return authHandler.RenderView(c, customerView)
}

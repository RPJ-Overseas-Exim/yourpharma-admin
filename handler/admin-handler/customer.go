package adminHandler

import (
	"RPJ-Overseas-Exim/yourpharma-admin/db/models"
	authHandler "RPJ-Overseas-Exim/yourpharma-admin/handler/auth-handler"
	"RPJ-Overseas-Exim/yourpharma-admin/pkg/types"
	adminView "RPJ-Overseas-Exim/yourpharma-admin/templ/admin-views"
	"bytes"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/Sydsvenskan/json2csv"
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
func (cs *customerService) GetCustomers(page int, limit int) ([]types.Customer, int, string, error) {
    var customersData []types.Customer
    var totalCustomers int
    result := cs.DB.Model(&models.Customer{}).Offset(page*limit).Limit(limit).Find(&customersData)
    if result.Error != nil {
        return customersData, 0, "", result.Error
    }

    // get the total customer counts
    result = cs.DB.Model(&models.Customer{}).Select("count(email) as Total").Find(&totalCustomers)
    if result.Error != nil {
        return customersData, totalCustomers, "", result.Error
    }

    // convert data to string
    var dataBuffer bytes.Buffer
    stringsByte, err := json.Marshal(customersData)
    if err == nil {
        json2csv.Convert(strings.NewReader(string(stringsByte)), &dataBuffer)
    }

    return customersData, totalCustomers, dataBuffer.String(), nil
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
    var deletedCustomer types.Customer
    result := cs.DB.Delete(&deletedCustomer, "id like ?", id)

    if result.Error != nil {
        log.Printf("\n\nFailed to delete the customer: %v", result.Error)
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
    customersData, totalCustomers, customersString, err := cs.GetCustomers(page, limit)
    if err != nil {
        log.Printf("Error in customers Data: %v", err)
    }

	customersView := adminView.Customers(customersData, totalCustomers, customersString, page, limit)
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

    customersData, totalCustomers, customersString, err := cs.GetCustomers(page, limit)
    if err != nil {
        log.Printf("Customers data is not fetched: %v", err)
    }

    customersView := adminView.Customers(customersData, totalCustomers, customersString, page, limit)
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

    customersData, totalCustomers, customersString, err := cs.GetCustomers(page, limit)
    customerView := adminView.Customers(customersData, totalCustomers, customersString, page, limit)
    return authHandler.RenderView(c, customerView)
}

func (cs *customerService) DeleteCustomer(c echo.Context) error {
    var err error
    id := c.Param("id")
    err = cs.DeleteCustomerDetails(id)
    if err != nil {
        log.Printf("\n\nfailed to delete the customer: %v", err)
    }

    page, err := strconv.Atoi(c.QueryParam("page"))
    if err!=nil {
        page = 0
    }
    limit, err := strconv.Atoi(c.QueryParam("limit"))
    if err != nil{
        limit = 10
    }

    customersData, totalCustomers, customersString, err := cs.GetCustomers(page, limit)
    customerView := adminView.Customers(customersData, totalCustomers, customersString, page, limit)
    return authHandler.RenderView(c, customerView)
}

package adminHandler

import (
	"RPJ-Overseas-Exim/yourpharma-admin/db/models"
	authHandler "RPJ-Overseas-Exim/yourpharma-admin/handler/auth-handler"
	"RPJ-Overseas-Exim/yourpharma-admin/pkg/types"
	"RPJ-Overseas-Exim/yourpharma-admin/pkg/utils"
	adminView "RPJ-Overseas-Exim/yourpharma-admin/templ/admin-views"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/Sydsvenskan/json2csv"
	"github.com/aidarkhanov/nanoid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type customerService struct {
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

func (cs *customerService) ImportCustomers(c echo.Context) error {
	csvFile, err := c.FormFile("csv-file")

	if err != nil {
		return err
	}

	src, err1 := csvFile.Open()

	if err1 != nil {
		return err1
	}

	fileScanner := bufio.NewScanner(src)
    cs.insertManyCustomers(fileScanner)

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 0
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}

    customersData, totalCustomers, customersString, err := cs.GetCustomers(page, limit)
    role := utils.GetRole(utils.GetAdmin(c))
    customerView := adminView.Customers(customersData, totalCustomers, customersString, page, limit, role)
    return authHandler.RenderView(c, customerView)
}

func (cs *customerService) insertManyCustomers(fileScanner *bufio.Scanner) error {
	fileScanner.Split(bufio.ScanLines)
	fileScanner.Scan()

	headers := strings.Split(fileScanner.Text(), ",")
	nameIndex := slices.IndexFunc(headers, func(heading string) bool { return strings.ToLower(heading) == "name" })
	emailIndex := slices.IndexFunc(headers, func(heading string) bool { return strings.ToLower(heading) == "email" })
	addressIndex := slices.IndexFunc(headers, func(heading string) bool { return strings.ToLower(heading) == "address" })
	numberIndex := slices.IndexFunc(headers, func(heading string) bool { return strings.ToLower(heading) == "number" })

	if nameIndex == -1 ||
		emailIndex == -1 {
		return fmt.Errorf("Name or Email not provided in the file")
	}

	var customers []*models.Customer
	var customer *models.Customer

	for fileScanner.Scan() {
		line := strings.Split(fileScanner.Text(), ",")
        name := line[nameIndex]
        email := line[emailIndex]

		if name == "" || email == "" {
			continue
		}

		matched, err := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]+$`, email)

		if err != nil && !matched {
			continue
		}

        var (
            address string
            number int
        )

		if addressIndex != -1 && line[addressIndex] != "" {
			address = line[addressIndex]
		}

		if numberIndex != -1 && line[numberIndex] != "" {
			if num, err := strconv.Atoi(line[numberIndex]); err == nil {
				number = num
			}
		}

		customer = models.NewCustomer(name, email, &number, address)
		customers = append(customers, customer)
	}

	err := cs.DB.Clauses(clause.OnConflict{
        Columns: []clause.Column{{Name:"email"}},
        DoUpdates: clause.AssignmentColumns([]string{"email"}),
    }).Create(&customers).Error
	return err
}

func (cs *customerService) AddCustomer(name, email string, number *int, address string) error {
	id := nanoid.New()
	customer := types.Customer{
		Id:      id,
		Name:    name,
		Email:   email,
		Number:  number,
		Address: address,
	}

	err := cs.DB.Create(&customer).Error

	if err != nil {
		return err
	}
	return nil
}

func (cs *customerService) UpdateCustomerDetails(id, name, email string, number *int, address string) error {
	result := cs.DB.Model(models.Customer{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":    name,
		"email":   email,
		"number":  number,
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

    role := utils.GetRole(utils.GetAdmin(c))
	customersView := adminView.Customers(customersData, totalCustomers, customersString, page, limit, role)
	return authHandler.RenderView(c, adminView.AdminIndex(
		"Customers",
		true,
		customersView,
        role,
	))
}

func (cs *customerService) CreateCustomer(c echo.Context) error {
	var err error
	num, err := strconv.Atoi(c.FormValue("number"))
	err = cs.AddCustomer(c.FormValue("name"), c.FormValue("email"), &num, c.FormValue("address"))

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 0
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}

    customersData, totalCustomers, customersString, err := cs.GetCustomers(page, limit)
    if err != nil {
        log.Printf("Customers data is not fetched: %v", err)
    }

    role := utils.GetRole(utils.GetAdmin(c))
    customersView := adminView.Customers(customersData, totalCustomers, customersString, page, limit, role)
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

	if err != nil {
		return err
	}

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 0
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}

    customersData, totalCustomers, customersString, err := cs.GetCustomers(page, limit)
    role := utils.GetRole(utils.GetAdmin(c))
    customerView := adminView.Customers(customersData, totalCustomers, customersString, page, limit, role)
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
	if err != nil {
		page = 0
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}

    customersData, totalCustomers, customersString, err := cs.GetCustomers(page, limit)
    role := utils.GetRole(utils.GetAdmin(c))
    customerView := adminView.Customers(customersData, totalCustomers, customersString, page, limit, role)
    return authHandler.RenderView(c, customerView)
}

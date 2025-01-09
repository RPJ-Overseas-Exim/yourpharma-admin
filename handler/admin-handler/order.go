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
	"mime/multipart"
	"regexp"
	"strconv"
	"strings"

	"github.com/Sydsvenskan/json2csv"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type orderService struct {
	DB *gorm.DB
}

func NewOrderService(db *gorm.DB) *orderService {
	return &orderService{DB: db}
}

// database functions ===================================================
func (ords *orderService) GetOrders(status string, limit, page int) ([]types.Order, int, string, error) {
	var ordersData []types.Order
	var totalOrders int
	var result *gorm.DB
	var countResult *gorm.DB
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
	dataQuery := ords.DB.Model(&models.Order{}).Select(selectStatement).Joins(joinStatement1).Joins(joinStatement2)
	countQuery := ords.DB.Model(&models.Order{}).Select("count(id) as Total")

	if len(status) != 0 && status != "all" {
		result = dataQuery.Where("status ilike ?", status).Offset(page * limit).Limit(limit).Scan(&ordersData)
		countResult = countQuery.Where("status ilike ?", status).Scan(&totalOrders)
	} else {
		result = dataQuery.Offset(page * limit).Limit(limit).Scan(&ordersData)
		countResult = countQuery.Scan(&totalOrders)
	}

	if result.Error != nil || countResult.Error != nil {
		return ordersData, totalOrders, "", result.Error
	}

	var dataBuffer bytes.Buffer
	stringsByte, err := json.Marshal(ordersData)
	if err == nil {
		json2csv.Convert(strings.NewReader(string(stringsByte)), &dataBuffer)
	}

	return ordersData, totalOrders, dataBuffer.String(), nil
}

func (ords *orderService) ImportOrders(c echo.Context) error {
	csvFile, err := c.FormFile("csv-file")
	if err != nil {
		return err
	}

	src, err1 := csvFile.Open()
	if err1 != nil {
		return err1
	}

    err = ords.insertManyOrders(src)
    if err!=nil {
        log.Printf("insert error: %v", err)
    }

    // handling the view part below
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 0
	}

	ordersData, totalOrders, ordersString, err := ords.GetOrders("", limit, page)
    if err != nil {
		log.Printf("Failed to get the order data: %v", err)
    }

    var productsData []types.Product
	result := ords.DB.Model(&models.Product{}).Select("name as Name").Scan(&productsData)
	if result.Error != nil {
		log.Printf("Failed to get the product data: %v", result.Error)
	}

	orderView := adminView.Orders(ordersData, "all", productsData, totalOrders, ordersString, limit, page)
	return authHandler.RenderView(c, orderView)

}

func (ords *orderService) getIndexMap(headings []string) map[string]int{
    var headingMap map[string]int = make(map[string]int, 0)
    for idx, heading := range(headings){
        headingMap[heading] = idx
    }

    return headingMap
}

func (ords *orderService) getProductsMap() map[string]models.Product{
    var productsMap map[string]models.Product = make(map[string]models.Product, 0)
    var products []models.Product
    result := ords.DB.Find(&products)
    if result.Error != nil {
        return productsMap
    }

    for _, product := range(products){
        productsMap[product.Name] = product
    }

    return productsMap
}

func (ords *orderService) getCustomersMap(customerNames []string) map[string]models.Customer{
    var customersMap map[string]models.Customer = make(map[string]models.Customer, 0)
    var customers []models.Customer
    result := ords.DB.Find(&customers, "name IN ?", customerNames)
    if result.Error != nil {
        return customersMap
    }

    for _, customer := range(customers){
        customersMap[customer.Email] = customer
    }

    return customersMap
}

func (ords *orderService) insertManyOrders(src multipart.File) error {
	fileScanner := bufio.NewScanner(src)
	fileScanner.Split(bufio.ScanLines)
	fileScanner.Scan()
	headers := strings.Split(fileScanner.Text(), ",")
    headingMap := ords.getIndexMap(headers)

    customerNameIndex := headingMap["name"] 
    customerEmailIndex := headingMap["email"]
    customerAddressIndex := headingMap["address"]
    customerNumberIndex := headingMap["number"]
    productNameIndex := headingMap["product"]
    productPriceIndex := headingMap["price"]
    productQtyIndex := headingMap["qty"]
    productStatusIndex := headingMap["status"] 

	if customerNameIndex == -1 ||
		customerEmailIndex == -1 ||
		customerNumberIndex == -1 ||
		customerAddressIndex == -1 ||
		productNameIndex == -1 ||
		productPriceIndex == -1 ||
		productQtyIndex == -1 ||
		productStatusIndex == -1 {
		return fmt.Errorf("name, email, address, product, price, or status is not provided in the file")
	}

    productsMap := ords.getProductsMap()

    var srcCopyString string
    var customerNames []string
    for fileScanner.Scan() {
        lineString := fileScanner.Text()
        line := strings.Split(lineString, ",")

        customerNames = append(customerNames, line[customerNameIndex])
        srcCopyString = srcCopyString + lineString
    }
    
    customersMap := ords.getCustomersMap(customerNames)

    fileScanner2 := bufio.NewScanner(bytes.NewBuffer([]byte(srcCopyString)))
	fileScanner2.Split(bufio.ScanLines)

	var orders []*models.Order

	for fileScanner2.Scan() {
		line := strings.Split(fileScanner2.Text(), ",")
        // log.Println("Line: ", line)

		customerEmail := line[customerEmailIndex]
		productName := strings.ToLower(line[productNameIndex])
		productStatus := strings.ToLower(line[productStatusIndex])

		productPrice, err := strconv.Atoi(line[productPriceIndex])
        if err != nil {
            log.Println("product price is not an int", productPriceIndex)
            continue
        }
		productQty, err := strconv.Atoi(line[productQtyIndex])
        if err != nil ||
			productPrice == 0 ||
			productQty == 0 ||
			customerEmail == "" ||
			productName == "" ||
			productStatus == "" {
            log.Println("product price, qty, name, status, customer email not provided")
			continue
		}

        var order *models.Order

		matched, err := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]+$`, customerEmail)
		if err != nil && !matched {
            log.Println("email is not valid")
			continue
		}

        product, ok := productsMap[productName]
        if !ok {
            log.Println("product is not present")
            continue
        }

        customer, ok := customersMap[customerEmail]
        if !ok {
            log.Println("customer is not present")
            continue
        }

        order = models.NewImportOrder(customer.Id, product.Id, productStatus, productQty, productPrice)
        orders = append(orders, order)

	}

    result := ords.DB.Create(&orders)
    if result.Error != nil {
        return result.Error
    }

	return nil
}


func (ords *orderService) AddOrderDetails(name, email, product string, number *int, quantity, price int, origin, address string) error {
	var customerData models.Customer
	var productData models.Product
	var result *gorm.DB

	result = ords.DB.Find(&customerData, "email = ?", email)
	if result.RowsAffected == 0 {
		newCustomer := models.NewCustomer(name, email, number, address)
		result = ords.DB.Create(newCustomer)
		customerData.Id = newCustomer.Id
	}

	result = ords.DB.Find(&productData, "name = ?", product)
	if result.RowsAffected == 0 {
		newProduct := models.NewProduct(product)
		result = ords.DB.Create(newProduct)
		productData.Id = newProduct.Id
	}

	newOrder := models.NewOrder(customerData.Id, productData.Id, origin, quantity, price)
	result = ords.DB.Create(newOrder)
	if result.Error != nil {
		log.Printf("error in order create: %v", result.Error)
	}

	return nil
}

func (ords *orderService) UpdateOrderDetails(id string) error {
	var order models.Order
	var status string

	result := ords.DB.Find(&order, "id = ?", id)
	utils.ErrorHandler(result.Error, "Failed to get the order details")

	if strings.ToLower(order.Status) == "active" {
		status = "paid"
	} else if strings.ToLower(order.Status) == "paid" {
		status = "shipped"
	} else if strings.ToLower(order.Status) == "shipped" {
		status = "delivered"
	}

	result = ords.DB.Model(&models.Order{}).Where("id like ?", id).Update("status", status)
	utils.ErrorHandler(result.Error, "Failed to update the order status")

	return nil
}

func (ords *orderService) DeleteOrderDetails(id string) error {
	var order types.Order
	result := ords.DB.Model(&models.Order{}).Where("id like ?", id).Delete(&order)
	utils.ErrorHandler(result.Error, "Failed to delete the order details")
	return nil
}

// routes functions ===================================================
func (ords *orderService) Orders(c echo.Context) error {
	var err error
	var productsData []types.Product
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 0
	}
	status := strings.ToLower(c.QueryParam("status"))

	ordersData, totalOrders, ordersString, err := ords.GetOrders(status, limit, page)
	if err != nil {
		log.Printf("Failed to get the order data: %v", err)
	}

	result := ords.DB.Model(&models.Product{}).Select("name as Name").Scan(&productsData)
	if result.Error != nil {
		log.Printf("Failed to get the product data: %v", result.Error)
	}
	ordersView := adminView.Orders(ordersData, status, productsData, totalOrders, ordersString, limit, page)

	return authHandler.RenderView(c, adminView.AdminIndex(
		"Orders",
		true,
		ordersView,
	))
}

func (ords *orderService) CreateOrder(c echo.Context) error {
	var err error
	var productsData []types.Product
	name := c.FormValue("name")
	email := c.FormValue("email")
	address := c.FormValue("address")
	product := c.FormValue("product")
	origin := "Dash"

	log.Printf("name: %v, email: %v, address: %v, product: %v", name, email, address, product)

	number, err := strconv.Atoi(c.FormValue("number"))
	utils.ErrorHandler(err, "Number is not provided")

	quantity, err := strconv.Atoi(c.FormValue("quantity"))
	utils.ErrorHandler(err, "Quantity is not provided")

	price, err := strconv.Atoi(c.FormValue("price"))
	utils.ErrorHandler(err, "Price is not provided")

	ords.AddOrderDetails(name, email, product, &number, quantity, price, origin, address)

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 0
	}

	ordersData, totalOrders, ordersString, err := ords.GetOrders("", limit, page)
	utils.ErrorHandler(err, "Failed to get the order data")
	result := ords.DB.Model(&models.Product{}).Select("name as Name").Scan(&productsData)
	if result.Error != nil {
		log.Printf("Failed to get the product data: %v", result.Error)
	}

	orderView := adminView.Orders(ordersData, "all", productsData, totalOrders, ordersString, limit, page)
	return authHandler.RenderView(c, orderView)
}

func (ords *orderService) UpdateOrder(c echo.Context) error {
	var err error
	var productsData []types.Product
	id := c.Param("id")

	ords.UpdateOrderDetails(id)

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 0
	}

	ordersData, totalOrders, ordersString, err := ords.GetOrders("", limit, page)
	utils.ErrorHandler(err, "Failed to get the order data")
	result := ords.DB.Model(&models.Product{}).Select("name as Name").Scan(&productsData)
	if result.Error != nil {
		log.Printf("Failed to get the product data: %v", result.Error)
	}

	orderView := adminView.Orders(ordersData, "all", productsData, totalOrders, ordersString, limit, page)
	return authHandler.RenderView(c, orderView)
}

func (ords *orderService) DeleteOrder(c echo.Context) error {
	var err error
	var productsData []types.Product
	id := c.Param("id")
	status := c.QueryParam("status")

	ords.DeleteOrderDetails(id)

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 0
	}

	ordersData, totalOrders, ordersString, err := ords.GetOrders(status, limit, page)
	utils.ErrorHandler(err, "Failed to get the order data")
	result := ords.DB.Model(&models.Product{}).Select("name as Name").Scan(&productsData)
	if result.Error != nil {
		log.Printf("Failed to get the product data: %v", result.Error)
	}

	orderView := adminView.Orders(ordersData, "all", productsData, totalOrders, ordersString, limit, page)
	return authHandler.RenderView(c, orderView)
}

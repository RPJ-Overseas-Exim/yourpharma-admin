package adminHandler

import (
	"RPJ-Overseas-Exim/yourpharma-admin/db/models"
	authHandler "RPJ-Overseas-Exim/yourpharma-admin/handler/auth-handler"
	"RPJ-Overseas-Exim/yourpharma-admin/pkg/types"
	"RPJ-Overseas-Exim/yourpharma-admin/pkg/utils"
	adminView "RPJ-Overseas-Exim/yourpharma-admin/templ/admin-views"
	"log"
	"strconv"
	"strings"

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
func (ords *orderService) GetOrders(status string, limit, page int) ([]types.Order, int, error) {
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
        result = dataQuery.Where("status like ?", status).Offset(page*limit).Limit(limit).Scan(&ordersData)
        countResult = countQuery.Where("status like ?", status).Scan(&totalOrders)
    }else{
        result = dataQuery.Offset(page*limit).Limit(limit).Scan(&ordersData)
        countResult = countQuery.Scan(&totalOrders)
    }

    if result.Error != nil || countResult.Error != nil {
        return ordersData, totalOrders, result.Error
    }

    return ordersData, totalOrders, nil
}

func (ords *orderService) AddOrderDetails(name, email, product string, number *int, quantity, price int, origin, address string) error {
    var customerData models.Customer
    var productData models.Product
    var result *gorm.DB 

    result = ords.DB.Find(&customerData, "email like ?", email)
    if result.RowsAffected == 0 {
        newCustomer := models.NewCustomer(name, email, number, address)
        result = ords.DB.Create(newCustomer)
        customerData.Id = newCustomer.Id
    }

    result = ords.DB.Find(&productData, "name like ?", product)
    if result.RowsAffected == 0 {
        newProduct := models.NewProduct(product)
        result = ords.DB.Create(newProduct)
        productData.Id = newProduct.Id
    }

    newOrder := models.NewOrder(customerData.Id, origin, productData.Id, quantity, price)
    result = ords.DB.Create(newOrder)
    return nil
}

func (ords *orderService) UpdateOrderDetails(id string) error {
    var order models.Order
    var status string

    result := ords.DB.Find(&order, "id like ?", id)
    utils.ErrorHandler(result.Error, "Failed to get the order details")

    if order.Status == "active" {
        status = "paid"
    }else if order.Status == "paid" {
        status = "shipped"
    }else if order.Status == "shipped" {
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

    ordersData, totalOrders, err := ords.GetOrders(status, limit, page)
    if err != nil {
        log.Printf("Failed to get the order data: %v", err)
    }

    result := ords.DB.Model(&models.Product{}).Select("name as Name").Scan(&productsData)
    if result.Error != nil {
        log.Printf("Failed to get the product data: %v", result.Error)
    }
	ordersView := adminView.Orders(ordersData, status, productsData, totalOrders, limit, page)

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
    if err!=nil{
        page = 0
    }

    ordersData, totalOrders, err := ords.GetOrders("", limit, page)
    utils.ErrorHandler(err, "Failed to get the order data")
    result := ords.DB.Model(&models.Product{}).Select("name as Name").Scan(&productsData)
    if result.Error != nil {
        log.Printf("Failed to get the product data: %v", result.Error)
    }

    orderView := adminView.Orders(ordersData, "All", productsData, totalOrders, limit, page)
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
    if err!=nil{
        page = 0
    }

    ordersData, totalOrders, err := ords.GetOrders("", limit, page)
    utils.ErrorHandler(err, "Failed to get the order data")
    result := ords.DB.Model(&models.Product{}).Select("name as Name").Scan(&productsData)
    if result.Error != nil {
        log.Printf("Failed to get the product data: %v", result.Error)
    }

    orderView := adminView.Orders(ordersData, "All", productsData, totalOrders, limit, page)
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
    if err!=nil{
        page = 0
    }

    ordersData, totalOrders, err := ords.GetOrders(status, limit, page)
    utils.ErrorHandler(err, "Failed to get the order data")
    result := ords.DB.Model(&models.Product{}).Select("name as Name").Scan(&productsData)
    if result.Error != nil {
        log.Printf("Failed to get the product data: %v", result.Error)
    }

    orderView := adminView.Orders(ordersData, "All", productsData, totalOrders, limit, page)
    return authHandler.RenderView(c, orderView)
}

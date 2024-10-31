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
func (ords *orderService) GetOrders(status string) ([]types.Order, error) {
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
        result = ords.DB.Model(&models.Order{}).Select(selectStatement).Joins(joinStatement1).Joins(joinStatement2).Where("orders.status like ?", status).Scan(&ordersData)
    }else{
        result = ords.DB.Model(&models.Order{}).Select(selectStatement).Joins(joinStatement1).Joins(joinStatement2).Scan(&ordersData)
    }

    log.Printf("Order result: %v", ordersData)

    if result.Error != nil {
        return ordersData, result.Error
    }

    return ordersData, nil
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
    status := strings.ToLower(c.QueryParam("status"))
    ordersData, err := ords.GetOrders(status)
    utils.ErrorHandler(err, "Failed to get the order data")

	ordersView := adminView.Orders(ordersData, status)
	var msgs []string

    log.Printf("status: %v", status)

	return authHandler.RenderView(c, adminView.AdminIndex(
		"Orders",
		true,
		msgs,
		msgs,
		ordersView,
	))
}

func (ords *orderService) CreateOrder(c echo.Context) error {
    var err error
    name := c.FormValue("name")
    email := c.FormValue("email")
    address := c.FormValue("address")
    product := c.FormValue("product")
    origin := "Dash"

    number, err := strconv.Atoi(c.FormValue("number"))
    utils.ErrorHandler(err, "Number is not provided")

    quantity, err := strconv.Atoi(c.FormValue("quantity"))
    utils.ErrorHandler(err, "Quantity is not provided")

    price, err := strconv.Atoi(c.FormValue("price"))
    utils.ErrorHandler(err, "Price is not provided")
    
    ords.AddOrderDetails(name, email, product, &number, quantity, price, origin, address)

    ordersData, err := ords.GetOrders("")
    utils.ErrorHandler(err, "Failed to get the order data")
    orderView := adminView.Orders(ordersData, "All")
    return authHandler.RenderView(c, orderView)
}

func (ords *orderService) UpdateOrder(c echo.Context) error {
    var err error
    id := c.Param("id")

    ords.UpdateOrderDetails(id)

    ordersData, err := ords.GetOrders("")
    utils.ErrorHandler(err, "Failed to get the order data")
    orderView := adminView.Orders(ordersData, "All")
    return authHandler.RenderView(c, orderView)
}

func (ords *orderService) DeleteOrder(c echo.Context) error {
    var err error
    id := c.Param("id")
    status := c.QueryParam("status")

    ords.DeleteOrderDetails(id)

    ordersData, err := ords.GetOrders(status)
    utils.ErrorHandler(err, "Failed to get the order data")
    orderView := adminView.Orders(ordersData, "All")
    return authHandler.RenderView(c, orderView)
}

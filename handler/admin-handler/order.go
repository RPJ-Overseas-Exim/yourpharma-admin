package adminHandler

import (
	"RPJ-Overseas-Exim/yourpharma-admin/db/models"
	authHandler "RPJ-Overseas-Exim/yourpharma-admin/handler/auth-handler"
	"RPJ-Overseas-Exim/yourpharma-admin/pkg/types"
	"RPJ-Overseas-Exim/yourpharma-admin/pkg/utils"
	adminView "RPJ-Overseas-Exim/yourpharma-admin/templ/admin-views"
	"log"
	"strconv"

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
func (ords *orderService) GetOrders() ([]types.Order, error) {
    var ordersData []types.Order

    result := ords.DB.Find(&ordersData)
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

func (ords *orderService) UpdataeOrderDetails(id string) error {
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
    ordersData, err := ords.GetOrders()
    utils.ErrorHandler(err, "Failed to get the order data")

	ordersView := adminView.Orders(ordersData)
	var msgs []string

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

    ordersData, err := ords.GetOrders()
    utils.ErrorHandler(err, "Failed to get the order data")
    orderView := adminView.Orders(ordersData)
    return authHandler.RenderView(c, orderView)
}

func (ords *orderService) UpdateOrder(c echo.Context) error {
    var err error
    ordersData, err := ords.GetOrders()
    utils.ErrorHandler(err, "Failed to get the order data")
    orderView := adminView.Orders(ordersData)
    return authHandler.RenderView(c, orderView)
}

func (ords *orderService) DeleteOrder(c echo.Context) error {
    var err error
    ordersData, err := ords.GetOrders()
    utils.ErrorHandler(err, "Failed to get the order data")
    orderView := adminView.Orders(ordersData)
    return authHandler.RenderView(c, orderView)
}

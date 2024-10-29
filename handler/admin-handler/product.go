package adminHandler

import (
	"RPJ-Overseas-Exim/yourpharma-admin/db/models"
	authHandler "RPJ-Overseas-Exim/yourpharma-admin/handler/auth-handler"
	"RPJ-Overseas-Exim/yourpharma-admin/pkg/types"
	"RPJ-Overseas-Exim/yourpharma-admin/pkg/utils"
	adminView "RPJ-Overseas-Exim/yourpharma-admin/templ/admin-views"
	"log"
	"strconv"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type productService struct{
    DB *gorm.DB
}

func NewProductService(db *gorm.DB) *productService {
    return &productService{DB:db}
}

func (ps *productService) ShowProductView(c echo.Context) templ.Component {
    var err error
    productsData, err := ps.GetProducts()

    if err != nil {
        log.Printf("Failed to get the products data: %v", err)
        return nil
    }

    productView := adminView.Products(productsData)

    return productView
}

// database functions --------------------------------------------------------
func (ps *productService) GetProducts() ([]types.Product, error) {
    var productsData []types.Product
    result := ps.DB.Model(&models.Product{}).Select("price_qties as Id, products.name as Name, price_qties.price as Price, price_qties.qty as Qty").Joins("left join price_qties on price_qties.product_id = products.id").Scan(&productsData)

    log.Printf("Product data: %v", productsData)

    if result.Error != nil {
        return productsData, result.Error
    }

    return productsData, nil
}

func (ps *productService) AddProductDetails(name string, qty, price int) error {
    newProduct := models.NewProduct(name)
    result := ps.DB.Create(&newProduct)
    if result.Error != nil {
        return result.Error
    }

    newPrice := models.NewPriceQty(newProduct.Id, qty, price)
    result = ps.DB.Create(&newPrice)
    if result.Error != nil {
        return result.Error
    }

    return nil
}

func (ps *productService) UpdateProductDetails(id, name string, price, qty int) error {
    newPriceQty := ps.DB.Model(&models.PriceQty{}).Where("id like ?", id).Updates(map[string]interface{}{
        "id": id,
        "price": price,
        "qty": qty,
    })

    utils.ErrorHandler(newPriceQty.Error, "Failed to update the price of the product")

    return nil
}

// routes functions --------------------------------------------------------
func (ps *productService) Products(c echo.Context) error {
    var msgs []string
    productView := ps.ShowProductView(c) 
   return authHandler.RenderView(c, adminView.AdminIndex(
        "Products",
        true,
        msgs,
        msgs,
        productView,
    ))

}

func (ps *productService) CreateProduct(c echo.Context) error {
    var err error
    quantity, err := strconv.Atoi(c.FormValue("quantity"))
    if err != nil {
        return err
    }
    price, err := strconv.Atoi(c.FormValue("price"))
    if err != nil {
        return err
    }

    ps.AddProductDetails(
        c.FormValue("name"),
        quantity,
        price,
    )

    productView := ps.ShowProductView(c)
    return authHandler.RenderView(c, productView)
}

func (ps *productService) UpdateProduct(c echo.Context) error {
    var err error 
    id := c.Param("id")
    name := c.FormValue("name")
    price, err := strconv.Atoi(c.FormValue("price"))
    utils.ErrorHandler(err, "Price is not provided")
    quantity, err := strconv.Atoi(c.FormValue("quantity"))
    utils.ErrorHandler(err, "Quantity is not provided")

    ps.UpdateProductDetails(id, name, price, quantity)

    productView := ps.ShowProductView(c)
    return authHandler.RenderView(c, productView)
}


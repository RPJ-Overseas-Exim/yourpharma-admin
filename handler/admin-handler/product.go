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

func (ps *productService) ProductView(c echo.Context) templ.Component {
    var err error
    productsData, err := ps.GetProducts()

    if err != nil {
        log.Printf("Failed to get the products data: %v", err)
        return nil
    }

    productView := adminView.Products(productsData)

    return productView
}

// database functions ========================================================
func (ps *productService) GetProducts() ([]types.Product, error) {
    var productsData []types.Product
    result := ps.DB.Model(&models.Product{}).Select("price_qties.id as Id, price_qties.product_id as PId, products.name as Name, price_qties.price as Price, price_qties.qty as Qty").Joins("inner join price_qties on price_qties.product_id = products.id").Find(&productsData)

    log.Printf("Product data: %v", productsData)

    if result.Error != nil {
        log.Printf("Product data error: %v", result.Error)
        return nil, result.Error
    }

    return productsData, nil
}

func (ps *productService) AddProductDetails(name string, qty, price int) error {
    var product models.Product
    var newPrice *models.PriceQty
    var result *gorm.DB

    result = ps.DB.Find(&product, "name like ?", name)
    utils.ErrorHandler(result.Error, "Failed to parse the product details")
    
    if  result.RowsAffected == 0 {
        newProduct := models.NewProduct(name)
        result = ps.DB.Create(&newProduct)
        if result.Error != nil {
            return result.Error
        }

        newPrice = models.NewPriceQty((*newProduct).Id, price, qty)
    }else{ 
        newPrice = models.NewPriceQty(product.Id, price, qty)
    }

    result = ps.DB.Create(&newPrice)
    if result.Error != nil {
        return result.Error
    }

    return nil
}

func (ps *productService) UpdateProductDetails(id, name string, price, qty int) error {
    var result *gorm.DB
    result = ps.DB.Model(&models.PriceQty{}).Where("id like ?", id).Updates(map[string]interface{}{
        "price": price,
        "qty": qty,
    })
    utils.ErrorHandler(result.Error, "Failed to update the price of the product")

    var updatedPrice models.PriceQty
    result = ps.DB.Find(&updatedPrice, "id like ?", id)
    utils.ErrorHandler(result.Error, "Failed to get the product details")

    result = ps.DB.Model(&models.Product{}).Where("id like ?", updatedPrice.ProductId).Update("name", name)
    utils.ErrorHandler(result.Error, "Failed to update the product name")

    return nil
}

func (ps *productService) DeleteProductDetails(id string) error {
    var result *gorm.DB
    var priceQty models.PriceQty
    var product models.Product

    result = ps.DB.Model(&models.PriceQty{}).Where("product_id like ?", id).Delete(&priceQty)
    utils.ErrorHandler(result.Error, "Failed to delete the prices of product")

    result = ps.DB.Model(&models.Product{}).Where("id like ?", id).Delete(&product)
    utils.ErrorHandler(result.Error, "Failed to delete the product")

   return nil 
}

func (ps *productService) DeletePriceDetails(id string) error {
    var price models.PriceQty
    result := ps.DB.Model(&models.PriceQty{}).Where("id like ?", id).Delete(&price)

    utils.ErrorHandler(result.Error, "Failed to delete the price details")
    
    return nil
}


// routes functions ========================================================
func (ps *productService) Products(c echo.Context) error {
    productView := ps.ProductView(c) 
   return authHandler.RenderView(c, adminView.AdminIndex(
        "Products",
        true,
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

    productView := ps.ProductView(c)
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

    productView := ps.ProductView(c)
    return authHandler.RenderView(c, productView)
}

func (ps *productService) DeletePrice(c echo.Context) error {
    id := c.Param("id")
    ps.DeletePriceDetails(id)
    productView := ps.ProductView(c)
    return authHandler.RenderView(c, productView)
}

func (ps *productService) DeleteProduct(c echo.Context) error {
    id := c.Param("id")
    ps.DeleteProductDetails(id)
    productView := ps.ProductView(c)
    return authHandler.RenderView(c, productView)
}

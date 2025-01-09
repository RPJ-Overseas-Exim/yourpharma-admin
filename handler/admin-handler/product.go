package adminHandler

import (
	"RPJ-Overseas-Exim/yourpharma-admin/db/models"
	authHandler "RPJ-Overseas-Exim/yourpharma-admin/handler/auth-handler"
	"RPJ-Overseas-Exim/yourpharma-admin/pkg/utils"
	adminView "RPJ-Overseas-Exim/yourpharma-admin/templ/admin-views"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/Sydsvenskan/json2csv"
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
    limit, err := strconv.Atoi(c.QueryParam("limit"))
    if err!=nil {
        limit = 10
    }
    page, err := strconv.Atoi(c.QueryParam("page"))
    if err!=nil {
        page=0
    }
    productsData, totalProducts, productsString, err := ps.GetProducts(limit, page)
    if err != nil {
        log.Printf("Failed to get the products data: %v", err)
        return nil
    }

    productView := adminView.Products(productsData, totalProducts, productsString, limit, page)
    return productView
}

// database functions ========================================================
func (ps *productService) GetProducts(limit, page int) ([]models.Product, int, string, error) {
    var priceData []models.PriceQty
    var productsData []models.Product

    result := ps.DB.Table("products").Preload("PriceQty").Find(&productsData)

    if result.Error != nil {
        log.Printf("Price data error: %v", result.Error)
        return productsData, 0, "", result.Error
    }

    result = ps.DB.Model(&models.PriceQty{}).Find(&priceData)
    if result.Error != nil {
        log.Printf("Price count error: %v", result.Error)
        return productsData, 0, "", result.Error
    }

    var dataBuffer bytes.Buffer
    stringsByte, err := json.Marshal(productsData)
    if err == nil {
        json2csv.Convert(strings.NewReader(string(stringsByte)), &dataBuffer)
    }

    return productsData, len(priceData), dataBuffer.String(), nil
}

func (ps *productService) AddProductDetails(name string) error {
        var result *gorm.DB
        var product models.Product
        name = strings.ToLower(name)

        result = ps.DB.Find(&product, "name = ?", name)

        if result.RowsAffected == 0 {
            newProduct := models.NewProduct(name)
            result = ps.DB.Create(&newProduct)
        }

        if result.Error != nil {
            return result.Error
        }

        return nil
}

func (ps *productService) AddPriceDetails(name string, qty, price int) error {
    var product models.Product
    var newPrice *models.PriceQty
    var result *gorm.DB

    result = ps.DB.Find(&product, "name = ?", name)
    utils.ErrorHandler(result.Error, "Failed to parse the product details")
    
    if  result.RowsAffected == 0 {
        newProduct := models.NewProduct(name)
        result = ps.DB.Create(&newProduct)
        if result.Error != nil {
            return result.Error
        }

        product.Id = newProduct.Id
    } 

    newPrice = models.NewPriceQty(product.Id, price, qty)
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

func (ps *productService) CreatePrice(c echo.Context) error {
    var err error
    quantity, err := strconv.Atoi(c.FormValue("quantity"))
    if err != nil {
        return err
    }
    price, err := strconv.Atoi(c.FormValue("price"))
    if err != nil {
        return err
    }

    ps.AddPriceDetails(
        c.FormValue("name"),
        quantity,
        price,
    )

    productView := ps.ProductView(c)
    return authHandler.RenderView(c, productView)
}

func (ps *productService) CreateProduct(c echo.Context) error {
    productName := c.FormValue("name") 
    if productName == "" {
        c.Response().WriteHeader(400)
        return errors.New("Product name should not be blank")
    }

    err := ps.AddProductDetails(productName)
    if err != nil {
        c.Response().WriteHeader(400)
        log.Printf("Failed to create the product: %v", err)
        return err
    }

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

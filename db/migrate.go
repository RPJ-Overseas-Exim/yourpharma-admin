package db

import (
	"RPJ-Overseas-Exim/yourpharma-admin/db/models"

	"gorm.io/gorm"
)


func migrate(db *gorm.DB){
    db.AutoMigrate(&models.PriceQty{}, &models.Product{}, &models.Customer{}, &models.Order{})
}

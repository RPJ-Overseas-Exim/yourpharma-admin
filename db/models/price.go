package models

type Price struct{
    Id          string      `json:"id" gorm:"primaryKey"`
    ProductId   Products    `json:"product_id" gorm:"foreignKey:ProductId"`
    Quantity    int         `json:"quantity" gorm:"not null"`
    Price       int         `json:"price" gorm:"not null"`
}

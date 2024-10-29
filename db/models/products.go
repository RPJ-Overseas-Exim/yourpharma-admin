package models

import "github.com/aidarkhanov/nanoid"

type Product struct {
    Id          string      `gorm:"primaryKey"`
    Name        string      `gorm:"uniqueIndex"`
    PriceQty    []PriceQty  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
    // Order       []Order     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}


func NewProduct(name string) *Product {
    id := nanoid.New()

    return &Product{
        Id:       id,
        Name:     name,
    }
}


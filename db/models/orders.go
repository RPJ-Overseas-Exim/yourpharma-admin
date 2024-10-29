package models

import (
	"time"

	"github.com/aidarkhanov/nanoid"
)

type Order struct {
    Id          string  
    CustomerId  string  `gorm:"uniqueIndex:idx_uniqueOrder"`
    ProductId   string  `gorm:"uniqueIndex:idx_uniqueOrder"`
    Quantity    int     `gorm:"uniqueIndex:idx_uniqueOrder"`
    Status      string  `gorm:"uniqueIndex:idx_uniqueOrder"`
    Origin      string
    Amount      int 
    CreatedAt,
    UpdatedAt   time.Time
}

func NewOrder(customerId, origin, productId string, quantity, amount int) *Order {
    id := nanoid.New()

    return &Order{
        Id: id,
        CustomerId: customerId,
        ProductId: productId,
        Quantity: quantity,
        Amount: amount,
        Status: "active",
        Origin: origin,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}

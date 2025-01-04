package models

import (
	"time"

	"github.com/aidarkhanov/nanoid"
)

type Order struct {
	Id,
    CustomerId string `gorm:"uniqueIndex:idx_uniqueOrder"`
    ProductId string `gorm:"uniqueIndex:idx_uniqueOrder"`
	Quantity int `gorm:"uniqueIndex:idx_uniqueOrder"`
	Amount int 
    Status string `gorm:"uniqueIndex:idx_uniqueOrder"`
    Origin string 
    CreatedAt,
    UpdatedAt time.Time
}

func NewOrder(customerId, productId, origin string, quantity, amount int) *Order{
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

func NewImportOrder(customerId, productId, status string, quantity, amount int) *Order{
    id := nanoid.New()

    return &Order{
        Id: id,
        CustomerId: customerId,
        ProductId: productId,
        Quantity: quantity,
        Amount: amount,
        Status: status,
        Origin: "Dash",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

}

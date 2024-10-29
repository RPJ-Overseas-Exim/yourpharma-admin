package models

import "github.com/aidarkhanov/nanoid"

type PriceQty struct {
    Id          string `gorm:"primaryKey"`
    ProductId   string
    Price, Qty  int
}

func NewPriceQty(prodId string, price, qty int) *PriceQty{
    id := nanoid.New()

    return &PriceQty{
        Id: id,
        ProductId: prodId,
        Price: price,
        Qty: qty,
    }
}

package models

import (
	"time"

	"github.com/aidarkhanov/nanoid"
)

type Product struct {
	Id       string
	Name     string     `gorm:"uniqueIndex"`
	PriceQty []PriceQty `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Order    []Order    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt,
	UpdatedAt time.Time
}

type PriceQty struct {
	Id,
	ProductId string
	Price, Qty int16
}

func NewPriceQty(prodId string, price, qty int) *PriceQty {
	id := nanoid.New()

	return &PriceQty{
		Id:        id,
		ProductId: prodId,
		Price:     int16(price),
		Qty:       int16(qty),
	}
}

func NewProduct(name string, price, quantity int) *Product {
	id := nanoid.New()
	priceQty := NewPriceQty(id, price, quantity)

	return &Product{
		Id:   id,
		Name: name,
		PriceQty: []PriceQty{
			*priceQty,
		},
	}
}

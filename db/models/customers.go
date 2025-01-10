package models

import (
	"time"

	"github.com/aidarkhanov/nanoid"
)

type Customer struct{
    Id,
    Name string 
    Email string `gorm:"uniqueIndex"`
    Number *int 
	Address *string
    Order []Order `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
    CreatedAt,
    UpdatedAt time.Time
}

func NewCustomer(name, email string, number *int, address string) *Customer{
    id := nanoid.New()

    return &Customer{
        Id:id,
        Name:name,
        Email:email,
        Number: number,
        Address: &address,
    }
}



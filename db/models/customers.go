package models

import (
	"github.com/aidarkhanov/nanoid"
	"gorm.io/gorm"
)

type Customer struct{
    gorm.Model
    Id          string  `gorm:"primaryKey"`
    Name        string  
    Email       string  `gorm:"uniqueIndex"`
    Number      *int  
    Address     string  
}

func NewCustomer(name, email string, number *int, address string) *Customer{
    id := nanoid.New()

    return &Customer{
        Id:id,
        Name:name,
        Email:email,
        Number: number,
        Address:address,
    }
}



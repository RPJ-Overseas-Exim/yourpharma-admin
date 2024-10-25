package models

import "time"

type Orders struct{
    Id          string      `json:"id" gorm:"primaryKey"`
    CustomerId  Customer    `json:"customer_id" gorm:"foreignKey:CustomerId"`
    PriceId     Price       `json:"price_id" gorm:"foreignKey:PriceId"`
    CreatedAt   time.Time   `json:"created_at"`
    Origin      string      `json:"origin"`
}


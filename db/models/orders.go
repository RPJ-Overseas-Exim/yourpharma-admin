package models

import "time"

type Orders struct{
    Id          uint        `json:"id" gorm:"primaryKey"`
    CustomerId  Customer    `json:"customer_id"`
    CreatedAt   time.Time   `json:"created_at" gorm:"unique"`
}


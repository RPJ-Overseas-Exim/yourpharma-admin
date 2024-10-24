package models

type Price struct{
    Id          int         `json:"id" gorm:"primaryKey"`
    Name        string      `json:"name" gorm:"not null"`
    Quantity    int         `json:"Quantity" gorm:"not null"`
}

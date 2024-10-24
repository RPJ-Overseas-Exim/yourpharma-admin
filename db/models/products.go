package models

type Products struct{
    Id          int    `json:"id" gorm:"primaryKey"`
    Name        string `json:"name" gorm:"not null"`
}



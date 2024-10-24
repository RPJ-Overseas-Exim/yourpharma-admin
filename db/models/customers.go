package models

type Customer struct{
    Id          uint    `json:"id" gorm:"primaryKey"`
    Name        string  `json:"name" gorm:"not null"`
    Email       string  `json:"email" gorm:"unique"`
    Number      string  `json:"number" gorm:"not null"`
    Address     string  `json:"address" gorm:"not null"`
}



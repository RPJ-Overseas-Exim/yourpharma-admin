package models

import (
	"errors"
	"log/slog"
	"time"

	"github.com/aidarkhanov/nanoid"
)

type Role int

const  (
    SuperAdminUser Role = iota 
    AdminUser
)

func (r Role) String () (string, error) {
    arr := [...]string{"super_admin", "admin"}
    if int(r)>=len(arr){
        return "", errors.New("No enum found")
    }
    return arr[r], nil
}

type Admin struct {
    Email string `gorm:"uniqueIndex"`
    Id,
    Password,
    Role string
    CreatedAt,
    UpdatedAt time.Time
}

func NewAdmin(email, password string, admin Role) *Admin{
    role,err := admin.String()
    if err!=nil{
        slog.Error("No role found for the given value", "role", role)
        return &Admin{}
    }

    return &Admin{
        Id: nanoid.New(),
        Email: email,
        Password: password,
        Role: role,
    }
}

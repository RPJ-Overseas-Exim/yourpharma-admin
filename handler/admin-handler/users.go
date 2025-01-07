package adminHandler

import (
	"RPJ-Overseas-Exim/yourpharma-admin/db/models"
	authHandler "RPJ-Overseas-Exim/yourpharma-admin/handler/auth-handler"
	adminView "RPJ-Overseas-Exim/yourpharma-admin/templ/admin-views"
	"log"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type userService struct{
    DB *gorm.DB
}

func NewUserService(db *gorm.DB) *userService{
    return &userService{
        DB: db,
    }
}

func (us *userService) GetUserPage(c echo.Context) error {
    var usersData []models.Admin

    result := us.DB.Find(&usersData)
    if result.Error != nil {
        log.Printf("error: %v", result.Error)
        return c.String(400, "Failed to get the users")
    }

    userView := adminView.Users(usersData)
    return authHandler.RenderView(c, adminView.AdminIndex(
        "Users",
        true,
        userView,
    ))
}

package adminHandler

import (
	"RPJ-Overseas-Exim/yourpharma-admin/handler/auth-handler"
	"RPJ-Overseas-Exim/yourpharma-admin/templ/admin-views"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type homeService struct{
    DB *gorm.DB
}

func NewHomeService(db *gorm.DB) *homeService {
    return &homeService{DB: db}
}

func (hs *homeService) Home(c echo.Context) error {
	homeView := adminView.Home()
	var msgs []string

	return authHandler.RenderView(c, adminView.AdminIndex(
		"Home",
		true,
		msgs,
		msgs,
		homeView,
	))
}


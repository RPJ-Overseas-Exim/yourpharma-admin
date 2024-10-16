package authHandlers

import (
	"RPJ-Overseas-Exim/yourpharma-admin/templ/authViews"

	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
    comp := authViews.LoginIndex("Login", authViews.Login())
    return comp.Render(c.Request().Context(), c.Response().Writer)
}

func Register(c echo.Context) error {
    comp := authViews.RegisterIndex("Register", authViews.Register())
    return comp.Render(c.Request().Context(), c.Response().Writer)
}

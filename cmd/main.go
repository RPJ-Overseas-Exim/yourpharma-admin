package main

import (
	// "RPJ-Overseas-Exim/yourpharma-admin/pkg/middleware"
	"RPJ-Overseas-Exim/yourpharma-admin/templ/authViews"
	"RPJ-Overseas-Exim/yourpharma-admin/templ/adminViews"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Data struct{
    Url string
}

func main(){
    e := echo.New()
    e.Use(middleware.Logger())
    e.Static("/static", "static")

    // middleware example
    // e.GET("/", func(c echo.Context) error {
    //     comp := authViews.LoginIndex("Login", authViews.Login())
    //     return comp.Render(c.Request().Context(), c.Response().Writer)
    // }, authMiddleware.AuthMiddleware)

    e.GET("/", func(c echo.Context) error {
        comp := authViews.LoginIndex("Login", authViews.Login())
        return comp.Render(c.Request().Context(), c.Response().Writer)
    })

    e.GET("/register", func(c echo.Context) error {
        comp := authViews.RegisterIndex("Register", authViews.Register())
        return comp.Render(c.Request().Context(), c.Response().Writer)
    })

    e.GET("/home", func(c echo.Context) error {
        comp := adminViews.HomeIndex("Home", adminViews.Home())
        return comp.Render(c.Request().Context(), c.Response().Writer)
    })

    e.Logger.Fatal(e.Start(":7000"))
}

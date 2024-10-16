package main

import (
	// "RPJ-Overseas-Exim/yourpharma-admin/pkg/middleware"
	"RPJ-Overseas-Exim/yourpharma-admin/handlers/adminHandlers"
	"RPJ-Overseas-Exim/yourpharma-admin/handlers/authHandlers"

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

    e.GET("/", authHandlers.Login)
    e.GET("/register", authHandlers.Register)
    e.GET("/home", adminHandlers.Home)
    e.GET("/customers", adminHandlers.Customers)

    e.Logger.Fatal(e.Start(":7000"))
}

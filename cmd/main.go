package main

import (
	"RPJ-Overseas-Exim/yourpharma-admin/handler"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Printf("Error in loading .env files %v", err)
        return 
    }
    log.Print("Env loaded successfully!")

	e := echo.New()

	e.Use(middleware.Logger())
	e.Static("/static", "static")

	handler.SetupAuthRoutes(e)
	handler.SetupAdminRoutes(e)

	e.Logger.Fatal(e.Start(":7000"))
}

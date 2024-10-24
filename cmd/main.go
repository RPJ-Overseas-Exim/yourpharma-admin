package main

import (
	"RPJ-Overseas-Exim/yourpharma-admin/db"
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

    db := db.ConnectDB()
    if db == nil {
        log.Print("Error in Database connection ")
        return 
    }
    log.Print("Database connected")

	e := echo.New()
    h := handlers.New(db)

	e.Use(middleware.Logger())
	e.Static("/static", "static")

	h.SetupAuthRoutes(e)
	h.SetupHomeRoutes(e)
	h.SetupCustomerRoutes(e)
	h.SetupOrderRoutes(e)

	e.Logger.Fatal(e.Start(":7000"))
}

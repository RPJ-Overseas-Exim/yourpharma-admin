package authMiddleware

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
   return func(c echo.Context) error {
        fmt.Println("auth middleware called")
        return next(c)
   }
}

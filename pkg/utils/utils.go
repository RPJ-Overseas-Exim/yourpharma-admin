package utils

import (
	"log/slog"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func GetAdmin(c echo.Context) *jwt.MapClaims{
    claims,ok := c.Get("admin").(jwt.MapClaims)
    if !ok{
        slog.Error("claims could not be generated")
    }
    return &claims
}

func GetRole(admin *jwt.MapClaims) string{
    return (*admin)["role"].(string)
}

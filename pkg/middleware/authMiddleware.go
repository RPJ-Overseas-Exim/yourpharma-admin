package authMiddleware

import (
	"RPJ-Overseas-Exim/yourpharma-admin/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
   return func(c echo.Context) error {
        token, err := c.Cookie("Authentication")
        if err != nil {
            return c.Redirect(http.StatusSeeOther, "/")
        }

        verifyErr := utils.VerifyToken(token.Value, []byte("secretKey"))
        // log.Printf("auth middleware called %v, err: %v", token, err)

        if verifyErr!=nil {
            return c.Redirect(http.StatusSeeOther, "/")
        }

        return next(c)
   }
}

func LoginMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func (c echo.Context) error {
        token , err := c.Cookie("Authentication")
        if err!=nil {
            return next(c)
        }

        verifyErr := utils.VerifyToken(token.Value, []byte("secretKey"))
        if verifyErr != nil {
            return next(c)
        }

        return c.Redirect(http.StatusSeeOther, "/home")
    }
}

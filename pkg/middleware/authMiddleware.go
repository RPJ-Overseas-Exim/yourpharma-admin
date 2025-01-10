package authMiddleware

import (
	"RPJ-Overseas-Exim/yourpharma-admin/pkg/utils"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
   return func(c echo.Context) error {

        err := godotenv.Load()
        if err != nil {
            c.Response().WriteHeader(500)
            return c.Redirect(http.StatusSeeOther, "/")
        }

        token, err := c.Cookie("Authentication")

        if err != nil {
            return c.Redirect(http.StatusSeeOther, "/")
        }

        verifyErr := utils.VerifyToken(token.Value, []byte(os.Getenv("JWT_SECRET")))
        // log.Printf("auth middleware called %v, err: %v", token, err)

        if verifyErr!=nil {
            return c.Redirect(http.StatusSeeOther, "/")
        }

        decoded, err := utils.DecodeToken(token.Value, []byte(os.Getenv("JWT_SECRET")))
        claims, ok := decoded.Claims.(jwt.MapClaims)
        if !ok{
            return c.Redirect(http.StatusSeeOther, "/")
        }

        c.Set("admin", claims)
        return next(c)
   }
}

func IsSuperAdmin(next echo.HandlerFunc) echo.HandlerFunc{
   return func(c echo.Context) error {
       adminRaw := c.Get("admin")
       admin, ok:= adminRaw.(jwt.MapClaims)

       if !ok {
           c.Response().Header().Set("HX-Redirect", "/")
           return c.String(401, "No admin found")
       }

       if role, ok1 := admin["role"]; ok1 && role=="super_admin"{
           return next(c)
       }

       c.Response().Header().Set("HX-Redirect", "/")
       return c.String(401, "Unauthorized")
   }
}

func LoginMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func (c echo.Context) error {
        token , err := c.Cookie("Authentication")

        if err!=nil {
            return next(c)
        }

        err = godotenv.Load()
        if err != nil {
            c.Response().WriteHeader(500)
            return next(c)
        }

        verifyErr := utils.VerifyToken(token.Value, []byte(os.Getenv("JWT_SECRET")))
        if verifyErr != nil {
            return next(c)
        }

        return c.Redirect(http.StatusSeeOther, "/home")
    }
}

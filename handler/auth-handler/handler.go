package authHandler

import (
	"RPJ-Overseas-Exim/yourpharma-admin/pkg/utils"
	authView "RPJ-Overseas-Exim/yourpharma-admin/templ/auth-views"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/a-h/templ"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type authService struct{
    DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *authService {
    return &authService{DB: db}
}

func (as *authService) LoginHandler(c echo.Context) error {
    var loginView templ.Component

	if c.Request().Method == "POST" {
        email := c.FormValue("email")
        password := c.FormValue("password")

        err := godotenv.Load()
        if err != nil {
            log.Printf("Failed to load the database url, %v", err)
            return nil
        }

        if email == os.Getenv("ADMIN_EMAIL") || password == os.Getenv("ADMIN_PASSWORD"){
            jwtCookie := new(http.Cookie)
            jwtCookie.Name = "Authentication"
            jwtCookie.Value = utils.CreateToken([]byte("secretKey"), "admin", os.Getenv("ADMIN_EMAIL"))
            jwtCookie.Expires = time.Now().Add(24 * time.Hour)

            c.SetCookie(jwtCookie)
            return c.Redirect(http.StatusSeeOther, "/home")
        }else{
   	        loginView = authView.Login("Email or password is incorrect")
            return RenderView(c, authView.LoginIndex(
                "Login",
                false,
                loginView,
            ))
        }
	}

   	loginView = authView.Login("")
	return RenderView(c, authView.LoginIndex(
		"Login",
		false,
		loginView,
	))
}

func RenderView(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

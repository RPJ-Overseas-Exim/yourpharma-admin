package authHandler

import (
	"RPJ-Overseas-Exim/yourpharma-admin/db/models"
	"RPJ-Overseas-Exim/yourpharma-admin/pkg/utils"
	authView "RPJ-Overseas-Exim/yourpharma-admin/templ/auth-views"
	"errors"
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
            c.Response().WriteHeader(500)
            return errors.New("Failed to load the database url")
        }

        var admin models.Admin
        as.DB.Find(&admin, "email = ? and password = ?", email, password)

        if admin.Email!=""{
            jwtCookie := new(http.Cookie)
            jwtCookie.Name = "Authentication"
            jwtCookie.Value = utils.CreateToken([]byte(os.Getenv("JWT_SECRET")), "admin", admin.Email, admin.Role)
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

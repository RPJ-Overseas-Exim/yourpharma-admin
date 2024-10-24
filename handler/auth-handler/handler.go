package authHandler

import (
	authView "RPJ-Overseas-Exim/yourpharma-admin/templ/auth-views"
	"encoding/json"
	"fmt"
	"log"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Admin struct{
    Email    string `json:"email"`
    Password string `json:"password"`
}

type authService struct{
    DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *authService {
    return &authService{DB: db}
}

func (as *authService) LoginHandler(c echo.Context) error {
   	loginView := authView.Login()
	var msgs []string

	if c.Request().Method == "POST" {
		// do the login procedures
        formData := make(map[string]interface{})
        err := json.NewDecoder(c.Request().Body).Decode(&formData)
        if err!=nil {
            log.Printf("The data is not right %v\n", err)
        }
        fmt.Printf("\nThe form data %v\n", formData)
	}

	return RenderView(c, authView.LoginIndex(
		"Login",
		false,
		msgs,
		msgs,
		loginView,
	))
}

func (as *authService) RegisterHandler(c echo.Context) error {
	registerView := authView.Register()
	var msgs []string

	if c.Request().Method == "POST" {
		// do the register procedure
	}

	return RenderView(c, authView.RegisterIndex(
		"Register",
		false,
		msgs,
		msgs,
		registerView,
	))
}

func RenderView(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

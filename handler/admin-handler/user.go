package adminHandler

import (
	"RPJ-Overseas-Exim/yourpharma-admin/db/models"
	authHandler "RPJ-Overseas-Exim/yourpharma-admin/handler/auth-handler"
	"RPJ-Overseas-Exim/yourpharma-admin/pkg/utils"
	adminView "RPJ-Overseas-Exim/yourpharma-admin/templ/admin-views"
	"errors"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type userServiceInt interface {
	getUsers() ([]models.Admin, error)
	createUser(email, password string) error
	updateUser(id, email, password string) error
	deleteUser(id string) error
}

type userService struct {
	dbConn *gorm.DB
}

type userHandler struct {
	us userServiceInt
}

func (uh *userHandler) GetUserPage(c echo.Context) error {
	usersData, err := uh.us.getUsers()

	if err != nil {
		c.Response().WriteHeader(500)
		return err
	}

	userView := adminView.Users(usersData)
    role := utils.GetRole(utils.GetAdmin(c))
	return authHandler.RenderView(c, adminView.AdminIndex(
		"Users",
		true,
		userView,
        role,
	))
}

func (uh *userHandler) HandleCreateUser(c echo.Context) error {
	error := uh.us.createUser(c.FormValue("email"), c.FormValue("password"))

	if error != nil {
		c.Response().WriteHeader(400)
		return error
	}

	return uh.getUserTable(c)
}

func (uh *userHandler) HandleUpdateUser(c echo.Context) error {
	error := uh.us.updateUser(c.Param("id"), c.FormValue("email"), c.FormValue("password"))
	if error != nil {
		c.Response().WriteHeader(400)
		return error
	}
	return uh.getUserTable(c)
}

func (uh *userHandler) HandleDeleteUser(c echo.Context) error {
	error := uh.us.deleteUser(c.Param("id"))
	if error != nil {
		c.Response().WriteHeader(400)
		return error
	}
	return uh.getUserTable(c)
}

func (uh *userHandler) getUserTable(c echo.Context) error {
	users, error := uh.us.getUsers()

	if error != nil {
		c.Response().WriteHeader(400)
		return error
	}

	userView := adminView.Users(users)
	return authHandler.RenderView(c, userView)
}

func (us *userService) getUsers() ([]models.Admin, error) {
	var admins []models.Admin
	err := us.dbConn.Find(&admins).Error
	return admins, err
}

func (us *userService) createUser(email, password string) error {
	if email == "" || password == "" {
		return errors.New("Please provide both email and password")
	}

	admin := models.NewAdmin(email, password, models.AdminUser)

	return us.dbConn.Create(admin).Error
}

func (us *userService) updateUser(id, email, password string) error {
	var admin models.Admin

	if id == "" || email == "" && password == "" {
		return errors.New("No updates to be done")
	}

	if email != "" {
		admin.Email = email
	}

	if password != "" {
		admin.Password = password
	}

	return us.dbConn.Model(&models.Admin{}).Where("id = ?", id).Updates(&admin).Error
}

func (us *userService) deleteUser(id string) error {
	if id == "" {
		return errors.New("No id provided to delete user")
	}

	var admin models.Admin

	admin.Email = id
	return us.dbConn.Where("id = ? and role <> ?", id, "super_admin").Delete(&admin).Error
}

func NewUserHandler(us userServiceInt) *userHandler {
	return &userHandler{
		us,
	}
}
func NewUserService(dbConn *gorm.DB) *userService {
	return &userService{
		dbConn,
	}
}

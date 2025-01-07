package adminHandler

import (
	"RPJ-Overseas-Exim/yourpharma-admin/db/models"
	authHandler "RPJ-Overseas-Exim/yourpharma-admin/handler/auth-handler"
	adminView "RPJ-Overseas-Exim/yourpharma-admin/templ/admin-views"
	"errors"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type userServiceInt interface {
    getUsers() ([]models.Admin, error)
    createUser(email, password string) error
    updateUser(email, password string) error
    deleteUser(email string) error
}

type userService struct{
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
    return authHandler.RenderView(c, adminView.AdminIndex(
        "Users",
        true,
        userView,
    ))
}

func (uh *userHandler) HandleCreateUser(c echo.Context) error{
    error := uh.us.createUser(c.Param("email"), c.Param("password"))

    if error!=nil{
        c.Response().WriteHeader(400)
        return error
    }

    return uh.GetUserPage(c)
}

func (uh *userHandler) HandleUpdateUser(c echo.Context) error{
    error := uh.us.updateUser(c.Param("email"), c.Param("password"))
    if error!=nil{
        c.Response().WriteHeader(400)
        return error
    }
    return uh.GetUserPage(c)
}

func (uh *userHandler) HandleDeleteUser(c echo.Context) error{
    error := uh.us.deleteUser(c.Param("email"))
    if error!=nil{
        c.Response().WriteHeader(400)
        return error
    }
    return uh.GetUserPage(c)
}

func (us *userService) getUsers() ([]models.Admin, error){
    var admins []models.Admin
    err := us.dbConn.Find(&admins).Error
    return admins,err
}

func (us *userService) createUser(email, password string) error {
    if email=="" || password==""{
        return errors.New("Please provide both email and password")
    }

    admin := models.NewAdmin(email, password, models.AdminUser)

    return us.dbConn.Create(admin).Error
}

func (us *userService) updateUser(email, password string) error {
    var admin models.Admin

    if email!=""{
        admin.Email = email
    }else if password!=""{
        admin.Password = password
    }else{
        return errors.New("No updates to be done")
    }

    return us.dbConn.Model(&models.Admin{}).Updates(&admin).Error
}

func (us *userService) deleteUser(email string) error{
    if email =="" {
        return errors.New("No email provided to delete user")
    }

    var admin models.Admin

    admin.Email = email
    return us.dbConn.Where("email = ?", email).Delete(&admin).Error
}

func NewUserHandler(us userServiceInt) *userHandler{
    return &userHandler{
        us,
    }
}
func NewUserService(dbConn *gorm.DB) *userService{
    return &userService{
        dbConn,
    }
}

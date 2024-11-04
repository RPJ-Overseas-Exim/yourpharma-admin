package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func CreateToken(secretKey []byte, username, email string) string {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
        "username": username,
        "email": email,
        "exp": time.Now().Add(time.Hour * 24).Unix(),
        "iat": time.Now().Unix(),
    })

    tokenString, err := token.SignedString(secretKey)
    if err!=nil {
        return ""
    }

    return tokenString
}

func VerifyToken(tokenString string, secretKey []byte) error {
    token , err := jwt.Parse(tokenString, func (token *jwt.Token) (interface{}, error){
        return secretKey, nil
    })

    if err!=nil {
        return err
    }

    if !token.Valid {
        return fmt.Errorf("Token is invalid")
    }
    
    return nil
}

func RemoveToken(c *echo.Context, tokenName string) {
    cookie := new(http.Cookie)
    cookie.Name = tokenName
    cookie.Value = ""
    (*c).SetCookie(cookie)
}

package configs

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
)

type JwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	Id    string `json:"id"`
	jwt.RegisteredClaims
}

// Configure middleware with the custom claims type

var Config = echojwt.Config{
	NewClaimsFunc: func(c echo.Context) jwt.Claims {
		return new(JwtCustomClaims)
	},
	SigningKey: []byte("secret"),
}

func GetClaims(c echo.Context) *JwtCustomClaims {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)
	return claims
}

func Restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)
	name := claims.Name
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

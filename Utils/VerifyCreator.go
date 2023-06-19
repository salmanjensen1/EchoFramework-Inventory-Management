package Utils

import (
	"RiseOfProduceManagement/configs"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func VerifyIfCreator(c echo.Context, idToCompare string) bool {
	claims := configs.GetClaims(c)
	id := claims.Id

	if id != idToCompare {
		return false
	}
	return true
}

func VerifyAdmin(c echo.Context) bool {
	user := c.Get("user").(*jwt.Token)
	fmt.Println(user)
	claims := configs.GetClaims(c)
	isAdmin := claims.Admin

	if !isAdmin {
		return false
	}
	return true
}

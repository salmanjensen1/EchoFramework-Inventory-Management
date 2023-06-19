package Middleware

import (
	"RiseOfProduceManagement/configs"
	"github.com/labstack/echo/v4"
)

func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//user := c.Get("user").(*jwt.Token)
		//claims, ok := user.Claims.(*jwtCustomClaims)
		//if !ok {
		//	fmt.Println("Claims is not ok >_<")
		//	return echo.ErrUnauthorized
		//}

		claims := configs.GetClaims(c)
		isAdmin := claims.Admin
		if isAdmin == false {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}

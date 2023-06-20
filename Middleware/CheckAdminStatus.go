package Middleware

import (
	"RiseOfProduceManagement/configs"
	"github.com/labstack/echo/v4"
)

func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		claims := configs.GetClaims(c)
		isAdmin := claims.Admin
		if isAdmin == false {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}

func IsNotAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		claims := configs.GetClaims(c)
		isAdmin := claims.Admin
		if isAdmin == true {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}

package Middleware

import (
	"RiseOfProduceManagement/configs"
	"fmt"
	"github.com/labstack/echo/v4"
	"time"
)

// ValidateToken validates the jwt token
func ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := configs.GetClaims(c)
		//fmt.Println("this is claims from token", claims)
		remainingTime := claims.ExpiresAt.Unix() - time.Now().Local().Unix()

		fmt.Println(remainingTime)

		if remainingTime <= 0 {
			return echo.ErrNotFound
		}

		return next(c)
	}
}

package main

import (
	"RiseOfProduceManagement/Routes"
	"RiseOfProduceManagement/configs"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

var productCollection *mongo.Collection = configs.GetCollection(configs.DB, "products")
var validate = validator.New()

// JwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examplesf

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*configs.JwtCustomClaims)
	name := claims.Name
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func isAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*configs.JwtCustomClaims)
		isAdmin := claims.Admin
		fmt.Println(claims.Name)
		if isAdmin == false {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}

// ValidateToken validates the jwt token
func validateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims, ok := user.Claims.(*configs.JwtCustomClaims)
		remainingTime := claims.ExpiresAt.Unix() - time.Now().Local().Unix()

		fmt.Println(remainingTime)

		if remainingTime <= 0 {
			return echo.ErrNotFound
		}

		if !ok {
			return echo.ErrUnauthorized
		}

		return next(c)
	}
}

func main() {
	e := echo.New()

	// Authentication routes
	Routes.NormalRoutes(e)
	Routes.AuthenticationRoutes(e)
	Routes.AdminRoutes(e)
	Routes.ProductRoutes(e)
	Routes.UserRoutes(e)

	// Restricted routes

	e.Logger.Fatal(e.Start(":1323"))
}

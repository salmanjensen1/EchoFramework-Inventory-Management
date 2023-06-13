package main

import (
	"RiseOfProduceManagement/Auth"
	"RiseOfProduceManagement/Response"
	"RiseOfProduceManagement/configs"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
	"time"
)

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func isAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*jwtCustomClaims)
		isAdmin := claims.Admin
		fmt.Println(claims.Name)
		if isAdmin == false {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}

// ValidateToken validates the jwt token
func ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims, ok := user.Claims.(*jwtCustomClaims)
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

func makeAdmin(c echo.Context) error {
	var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

	isToMakeAdminString := c.FormValue("id")
	isToMakeAdmin, _ := strconv.Atoi(isToMakeAdminString)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//get the ID for which the data is to be updated against
	updateUser := bson.M{
		"isadmin": true,
	}

	//insert the updated student info against the received id in database
	result, err := userCollection.UpdateOne(ctx, bson.M{"id": isToMakeAdmin}, bson.M{"$set": updateUser})

	if err != nil || result.MatchedCount != 1 {
		return c.String(500, "Update failed")
	}

	return c.JSON(200, Response.SystemResponse{200, "You are admin",
		&echo.Map{"data": result}})
}

func main() {
	e := echo.New()

	// Login route
	e.POST("/login", Auth.Login)
	//Register route
	e.POST("/register", Auth.Register)
	// Unauthenticated route
	e.GET("/", accessible)

	// Restricted group
	r := e.Group("/restricted")

	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte("secret"),
	}
	r.Use(echojwt.WithConfig(config))
	a := r.Group("/forAdmin", ValidateToken, isAdmin)
	a.GET("/", restricted)
	a.POST("/make-admin", makeAdmin)
	r.GET("/forUser", restricted, ValidateToken)
	e.Logger.Fatal(e.Start(":1323"))
}

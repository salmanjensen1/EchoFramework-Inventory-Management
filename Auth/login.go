package Auth

import (
	"RiseOfProduceManagement/Model"
	"RiseOfProduceManagement/Response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	Id    string `json:"id"`
	jwt.RegisteredClaims
}

func Login(c echo.Context) error {
	ctx, _ := context.WithTimeout(context.Background(), 50*time.Second)
	//defer cancel()

	var user Model.User

	reqUsername := c.FormValue("username")
	reqPassword := c.FormValue("password")

	// Throws unauthorized error
	err := userCollection.FindOne(ctx, bson.M{"username": reqUsername}).Decode(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response.SystemResponse{Status: http.StatusInternalServerError, Message: "username doesn't exist",
			Data: &echo.Map{"data": err.Error()}})
	}

	if user.DeleteStatus == true {
		return c.String(404, "Your account has been deleted")
	}

	if err1 := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqPassword)); err1 != nil {
		return c.JSON(http.StatusInternalServerError, Response.SystemResponse{Status: http.StatusInternalServerError, Message: "Invalid Password",
			Data: &echo.Map{"data": err1.Error()}})
	}
	userIDString := user.ID.Hex()
	// Set custom claims
	claims := &jwtCustomClaims{
		user.Name,
		user.IsAdmin,
		userIDString,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

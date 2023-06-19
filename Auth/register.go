package Auth

import (
	"RiseOfProduceManagement/Model"
	"RiseOfProduceManagement/Response"
	"RiseOfProduceManagement/configs"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

func Register(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user Model.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, Response.SystemResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.JSON(http.StatusBadRequest, Response.SystemResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": validationErr.Error()}})
	}

	// Check in your db if the user already exists or not
	count, err := userCollection.CountDocuments(ctx, bson.M{"username": user.Username})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response.SystemResponse{Status: http.StatusInternalServerError, Message: "Error looking up for document",
			Data: &echo.Map{"data": err.Error()}})
	}

	if count > 0 {
		return c.String(http.StatusInternalServerError, "this email or username already exists")

	}

	count, err = userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response.SystemResponse{Status: http.StatusInternalServerError, Message: "Error looking up for document (Email)",
			Data: &echo.Map{"data": err.Error()}})
	}
	fmt.Println(count)
	if count > 0 {
		return c.String(http.StatusInternalServerError, "this email or username already exists")

	}

	passwordHashed := HashPassword(user.Password)
	user.ID = primitive.NewObjectID()

	newUser := Model.User{
		ID:             user.ID,
		Name:           user.Name,
		Username:       user.Username,
		Password:       passwordHashed,
		Email:          user.Email,
		Phone:          user.Phone,
		Address:        user.Address,
		IsAdmin:        false,
		DeleteStatus:   false,
		AccountBalance: 500,
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response.SystemResponse{Status: http.StatusInternalServerError, Message: "Couldn't insert data", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusCreated, Response.SystemResponse{Status: http.StatusCreated, Message: "success", Data: &echo.Map{"data": result}})

}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(string(bytes))

	return string(bytes)
}

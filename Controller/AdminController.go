package Controller

import (
	"RiseOfProduceManagement/Response"
	"RiseOfProduceManagement/configs"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"strconv"
	"time"
)

func MakeAdmin(c echo.Context) error {
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

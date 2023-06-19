package Controller

import (
	"RiseOfProduceManagement/Model"
	"RiseOfProduceManagement/Response"
	"RiseOfProduceManagement/configs"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"time"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

func MakeAdmin(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	isToMakeAdminString := c.FormValue("id")
	isToMakeAdmin, _ := primitive.ObjectIDFromHex(isToMakeAdminString)

	var user Model.User

	filterUser := bson.M{"_id": isToMakeAdmin}
	err := userCollection.FindOne(ctx, filterUser).Decode(&user)

	//if the user is not found
	if err != nil {
		return c.String(500, "User not found")
	}

	//if the user is already an admin
	if user.IsAdmin {
		return c.String(500, user.Name+" is already an admin")
	}

	//update the user to make him an admin
	updateUser := bson.M{"$set": bson.M{"isadmin": true}}

	//insert the updated student info against the received id in database
	result, err := userCollection.UpdateOne(ctx, filterUser, updateUser)
	result.UpsertedID = user.ID
	if err != nil || result.MatchedCount != 1 {
		return c.String(500, "Update failed")
	}

	return c.JSON(200, Response.SystemResponse{200, "You are admin",
		&echo.Map{"data": result}})
}

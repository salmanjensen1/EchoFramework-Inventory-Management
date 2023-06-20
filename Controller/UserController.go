package Controller

import (
	"RiseOfProduceManagement/Auth"
	"RiseOfProduceManagement/Model"
	"RiseOfProduceManagement/Response"
	"RiseOfProduceManagement/Utils"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func ViewProfile(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID := c.Param("userID")
	userIDObject, _ := primitive.ObjectIDFromHex(userID)

	filter := bson.M{"_id": userIDObject}
	var user Model.User
	err1 := userCollection.FindOne(ctx, filter).Decode(&user)

	if user.DeleteStatus == true {
		return c.String(404, "User is deleted")
	}

	if err1 != nil {
		return c.JSON(http.StatusBadRequest, Response.SystemResponse{400, "Couldn't find the user from DB",
			&echo.Map{"data": err1.Error()}})
	}

	mapOfUsers := make(map[string]userDetails)

	mapOfUsers[user.Name] = userDetails{
		Name:           user.Name,
		Username:       user.Username,
		Email:          user.Email,
		Phone:          user.Phone,
		Address:        user.Address,
		AccountBalance: user.AccountBalance,
	}
	// Print the empty ObjectID

	return c.JSON(200, Response.SystemResponse{200, "User profile info: ", &echo.Map{"data": mapOfUsers}})
}

func UpdateUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var reqUserData, user Model.User

	//get the data that is to be updated
	if err := c.Bind(&reqUserData); err != nil {
		return c.String(http.StatusBadRequest, "the request data is not acceptable")
	}

	if validationErr := validate.Struct(&reqUserData); validationErr != nil {
		return c.String(http.StatusBadRequest, "Data Validation failed")
	}

	if reqUserData.DeleteStatus == true {
		return c.String(http.StatusBadRequest, "User is deleted")
	}

	//get the ID for which the data is to be updated against
	updateID := c.Param("userID")
	updateIDObject, _ := primitive.ObjectIDFromHex(updateID)

	//check if the user exists in the database
	filter := bson.M{"_id": updateIDObject}
	err1 := userCollection.FindOne(ctx, filter).Decode(&user)

	if err1 != nil {
		return c.JSON(http.StatusBadRequest, Response.SystemResponse{400, "Couldn't find the user",
			&echo.Map{"data": err1.Error()}})
	}
	//*****************************************************
	creator := Utils.VerifyIfCreator(c, updateID)
	if !creator {
		return c.String(http.StatusForbidden, "You are not authorized to update this user")
	}
	newUser := Model.User{
		ID:             user.ID,
		Name:           reqUserData.Name,
		Username:       reqUserData.Username,
		Password:       Auth.HashPassword(reqUserData.Password),
		Email:          reqUserData.Email,
		Phone:          reqUserData.Phone,
		Address:        reqUserData.Address,
		IsAdmin:        user.IsAdmin,
		DeleteStatus:   user.DeleteStatus,
		AccountBalance: user.AccountBalance,
	}
	//insert the updated user info against the received id in database
	result, err1 := userCollection.UpdateOne(ctx, bson.M{"_id": updateIDObject}, bson.M{"$set": newUser})

	if err1 != nil || result.MatchedCount != 1 {
		return c.String(500, "Update failed")
	}

	return c.JSON(200, Response.SystemResponse{200, "Updating user Entry Complete",
		&echo.Map{"data": result}})

}

func DeleteUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	deleteIdString := c.Param("userID")
	deleteIdObject, _ := primitive.ObjectIDFromHex(deleteIdString)

	var user Model.User
	filter := bson.M{"_id": deleteIdObject}
	err1 := userCollection.FindOne(ctx, filter).Decode(&user)

	if user.DeleteStatus == true {
		return c.String(http.StatusBadRequest, "User is already deleted")
	}

	if err1 != nil {
		return c.JSON(http.StatusBadRequest, Response.SystemResponse{400, "Couldn't find the user",
			&echo.Map{"data": err1.Error()}})
	}

	creator := Utils.VerifyIfCreator(c, deleteIdString)
	isAdmin := Utils.VerifyAdmin(c)

	fmt.Println(creator, isAdmin)

	if !creator && !isAdmin {
		return c.String(http.StatusForbidden, "You are not authorized to delete this user")
	}

	//result, err := userCollection.DeleteOne(ctx, bson.M{"_id": deleteIdObject})

	//if err != nil || result.DeletedCount < 1 {
	//	return c.String(500, "Couldn't delete the user Entry")
	//}

	//safe delete a user
	result, err := userCollection.UpdateOne(ctx, bson.M{"_id": deleteIdObject}, bson.M{"$set": bson.M{"deletestatus": true}})
	if err != nil || result.MatchedCount < 1 {
		return c.String(500, "Couldn't delete the user Entry")
	}

	//cascade delete all the products of the user
	DeleteAllProductsOfTheSeller(c, deleteIdString)
	return c.String(200, "deleted user: "+user.Name+" successfully")
}

func DeleteAllProductsOfTheSeller(c echo.Context, sellerID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"sellerid": sellerID}

	products, err := productCollection.Find(ctx, filter)

	if products.RemainingBatchLength() < 1 {
		return c.String(http.StatusNotFound, "No Products found")
	}
	mapOfProducts := make(map[string]productDetails)
	for products.Next(context.TODO()) {
		var product Model.Product
		if err := products.Decode(&product); err != nil {
			// Handle the error
			return c.String(400, "Error decoding products from database")
		}
		fmt.Println(product)
		product.DeleteStatus = true
		productIDString := product.ID.Hex()

		// Do something with the user document
		mapOfProducts[productIDString] = productDetails{
			Name:     product.ProductName,
			Quantity: product.Quantity,
			Price:    product.Price,
		}
		result, err := productCollection.UpdateOne(ctx, bson.M{"_id": product.ID}, bson.M{"$set": product})
		if err != nil || result.MatchedCount < 1 {
			return c.String(500, "Couldn't update the product Entry")
		}

	}

	// Close the cursor
	products.Close(context.TODO())

	if err != nil {
		return c.JSON(http.StatusNotFound, Response.SystemResponse{400, "Error querying products from database",
			&echo.Map{"data": err.Error()}})
	}
	return c.JSON(200, Response.SystemResponse{200, "Deleted products", &echo.Map{"data": mapOfProducts}})
}

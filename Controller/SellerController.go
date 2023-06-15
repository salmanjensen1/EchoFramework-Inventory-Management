package Controller

import (
	"RiseOfProduceManagement/Model"
	"RiseOfProduceManagement/Response"
	"context"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strconv"
	"time"
)

func CreateSeller(c echo.Context) error {
	//close this function if it takes more than 10 seconds time
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer is called whether the parent/surrounding function is finished or not
	defer cancel()

	sellerID := c.Param("sellerID")

	//create an instance of the product model (expected incoming request)
	var product Model.Product
	//bind the incoming json data
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, Response.SystemResponse{400, "The Input data doesn't match the input fields",
			&echo.Map{"data": err.Error()}})
	}

	//check if all the Validation constraint are met specified in the Predefined Response Struct

	product.SellerID = sellerID

	if validationError := validate.Struct(&product); validationError != nil {
		return c.JSON(http.StatusBadRequest, Response.SystemResponse{400, "Data Validation failed. Input data in all required fields",
			&echo.Map{"data": validationError.Error()}})
	}
	result, err := productCollection.InsertOne(ctx, product)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response.SystemResponse{500, "Error inserting in Database",
			&echo.Map{"data": err.Error()}})
	}
	return c.JSON(200, Response.SystemResponse{200, "New Product Entry Complete",
		&echo.Map{"data": result}})
}

func UpdateSeller(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var product Model.Product

	//get the data that is to be updated
	if err := c.Bind(&product); err != nil {
		return c.String(http.StatusBadRequest, "the request data is not acceptable")
	}

	if validationErr := validate.Struct(&product); validationErr != nil {
		return c.String(http.StatusBadRequest, "Data Validation failed")
	}

	//get the ID for which the data is to be updated against
	updateIDString := c.Param("id")
	updateID, _ := strconv.Atoi(updateIDString) //convert the ID received from params from string to int

	//insert the updated product info against the received id in database
	result, err := productCollection.UpdateOne(ctx, bson.M{"id": updateID}, bson.M{"$set": product})

	if err != nil || result.MatchedCount != 1 {
		return c.String(500, "Update failed")
	}

	return c.JSON(200, Response.SystemResponse{200, "Updating Product Entry Complete",
		&echo.Map{"data": result}})

}

func DeleteSeller(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	deleteIdString := c.Param("id")
	deleteId, _ := strconv.Atoi(deleteIdString)

	result, err := productCollection.DeleteOne(ctx, bson.M{"id": deleteId})

	if err != nil || result.DeletedCount < 1 {
		return c.String(500, "Couldn't delete the Product Entry")
	}

	return c.String(200, "deleted Product successfully")
}

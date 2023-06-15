package Controller

import (
	"RiseOfProduceManagement/Model"
	"RiseOfProduceManagement/Response"
	"RiseOfProduceManagement/configs"
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
	"time"
)

var productCollection *mongo.Collection = configs.GetCollection(configs.DB, "products")
var validate = validator.New()

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	ID    string `json:"id"`
	jwt.RegisteredClaims
}

type productDetails struct {
	Name     string
	Quantity int
	Price    float64
}

func CreateProduct(c echo.Context) error {
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
	product.ID = primitive.NewObjectID()

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

func GetProduct(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	productID := c.Param("productID")
	productIDObject, _ := primitive.ObjectIDFromHex(productID)

	filter := bson.M{"_id": productIDObject}
	var product Model.Product
	err1 := productCollection.FindOne(ctx, filter).Decode(&product)

	if err1 != nil {
		return c.JSON(http.StatusBadRequest, Response.SystemResponse{400, "Couldn't find the product",
			&echo.Map{"data": err1.Error()}})
	}

	newProduct := Model.Product{
		ProductName: product.ProductName,
		Quantity:    product.Quantity,
	}
	// Print the empty ObjectID

	return c.JSON(200, Response.SystemResponse{200, "Product found", &echo.Map{"data": newProduct}})
}

func SearchProduct(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	productName := c.Param("productName")
	fmt.Println(productName)
	filter := bson.M{"productname": productName}

	products, err1 := productCollection.Find(ctx, filter)
	mapOfProducts := make(map[string]productDetails)

	fmt.Println(products)
	for products.Next(context.TODO()) {
		var product Model.Product
		if err := products.Decode(&product); err != nil {
			// Handle the error
		}
		productIDString := product.ID.Hex()

		// Do something with the user document

		mapOfProducts[productIDString] = productDetails{
			Name:     product.ProductName,
			Quantity: product.Quantity,
			Price:    product.Price,
		}
	}

	// Close the cursor
	products.Close(context.TODO())
	//condition to check if the product is not found and at least one product is found
	if err1 != nil || products == nil {
		return c.JSON(http.StatusBadRequest, Response.SystemResponse{400, "Couldn't find the product",
			&echo.Map{"data": err1.Error()}})
	}
	return c.JSON(200, Response.SystemResponse{200, "Product found", &echo.Map{"data": mapOfProducts}})
}

func GetAllProductsOfASeller(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sellerID := c.Param("sellerID")

	fmt.Println(sellerID)

	filter := bson.M{"sellerid": sellerID}

	products, err := productCollection.Find(ctx, filter)

	mapOfProducts := make(map[string]productDetails)
	for products.Next(context.TODO()) {
		var product Model.Product
		if err := products.Decode(&product); err != nil {
			// Handle the error
		}
		productIDString := product.ID.Hex()

		// Do something with the user document

		mapOfProducts[productIDString] = productDetails{
			Name:     product.ProductName,
			Quantity: product.Quantity,
			Price:    product.Price,
		}
	}

	// Close the cursor
	products.Close(context.TODO())
	if err != nil {
		return c.JSON(http.StatusNotFound, Response.SystemResponse{400, "No Products found",
			&echo.Map{"data": err.Error()}})
	}
	return c.JSON(200, Response.SystemResponse{200, "Product found", &echo.Map{"data": mapOfProducts}})
}

func UpdateProduct(c echo.Context) error {
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
	updateID := c.Param("productID")
	updateIDObject, _ := primitive.ObjectIDFromHex(updateID)

	err := productCollection.FindOne(ctx, bson.M{"_id": updateIDObject}).Decode(&product)
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	id := claims.ID

	if id != product.SellerID {
		fmt.Println(err)
		return c.String(http.StatusForbidden, "You are not authorized to update this product")
	}

	//insert the updated product info against the received id in database
	result, err1 := productCollection.UpdateOne(ctx, bson.M{"_id": updateIDObject}, bson.M{"$set": product})

	if err1 != nil || result.MatchedCount != 1 {
		return c.String(500, "Update failed")
	}

	return c.JSON(200, Response.SystemResponse{200, "Updating Product Entry Complete",
		&echo.Map{"data": result}})

}

func DeleteProduct(c echo.Context) error {
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

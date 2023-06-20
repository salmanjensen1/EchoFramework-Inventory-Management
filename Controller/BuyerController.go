package Controller

import (
	"RiseOfProduceManagement/Model"
	"RiseOfProduceManagement/Response"
	"RiseOfProduceManagement/configs"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"strconv"
	"time"
)

func BuyProduct(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer is called whether the parent/surrounding function is finished or not
	defer cancel()

	productID := c.QueryParam("productID")
	productQtyString := c.QueryParam("productQty")

	//convert the productID string to an int
	productQty, _ := strconv.Atoi(productQtyString)

	productIDObject, _ := primitive.ObjectIDFromHex(productID)
	//create an instance of the product model (expected incoming request)
	var product Model.Product
	//find the product with the given ID
	filterProduct := bson.M{"_id": productIDObject}
	err := productCollection.FindOne(ctx, filterProduct).Decode(&product)

	if err != nil {
		return c.JSON(404, Response.SystemResponse{404, "Product not found", &echo.Map{"data": err.Error()}})

	}

	claims := configs.GetClaims(c)
	fmt.Println(claims)
	var buyer, seller Model.User

	//get buyer ID from token
	buyerID := claims.Id
	buyerIDObject, err := primitive.ObjectIDFromHex(buyerID)

	if err != nil {
		c.String(500, "Error converting buyer ID to ObjectID")
	}

	filterBuyer := bson.M{"_id": buyerIDObject}
	//find the buyer with the given ID
	fmt.Println(buyerIDObject)
	err1 := userCollection.FindOne(ctx, bson.M{"_id": buyerIDObject}).Decode(&buyer)
	fmt.Println(buyer)
	if err1 != nil {
		return c.JSON(500, Response.SystemResponse{500, "Buyer ID not found", &echo.Map{"data": err1.Error()}})
	}

	//find seller info from product's seller ID field
	sellerID := product.SellerID
	sellerIDObject, _ := primitive.ObjectIDFromHex(sellerID)
	filterSeller := bson.M{"_id": sellerIDObject}
	//find the buyer with the given ID
	err10 := userCollection.FindOne(ctx, filterSeller).Decode(&seller)

	if err10 != nil {
		return c.JSON(500, Response.SystemResponse{500, "Seller ID" + sellerID + "not found", &echo.Map{"data": err10.Error()}})
	}

	//input validation
	if productQty < 1 {
		return c.String(500, "Invalid quantity")
	}

	//calculate remaining quantity
	remainingQty := product.Quantity - productQty
	//calculate amount
	amount := product.Price * float64(productQty)
	//deduct the quantity of the product
	if remainingQty < 1 {
		return c.String(500, "The quantity requested is more than the quantity available")
	}
	//check if the user has sufficient fund to buy a product
	if amount > buyer.AccountBalance {
		return c.String(500, "Insufficient Funds")
	}
	//deduct the amount from the buyer who bought the product
	buyer.AccountBalance = buyer.AccountBalance - amount
	seller.AccountBalance = seller.AccountBalance + amount
	product.Quantity = product.Quantity - productQty

	//update the user balance in the database
	result, err2 := userCollection.UpdateOne(ctx, filterBuyer, bson.M{"$set": buyer})
	if err2 != nil || result.MatchedCount != 1 {
		return c.String(500, "Error updating buyer account balance")
	}

	//update the product quantity in the database
	result1, err3 := productCollection.UpdateOne(ctx, filterProduct, bson.M{"$set": product})
	if err3 != nil || result1.MatchedCount != 1 {
		return c.String(500, "Update of Product Qty failed")
	}

	result2, err3 := userCollection.UpdateOne(ctx, filterSeller, bson.M{"$set": seller})
	if err3 != nil || result2.MatchedCount != 1 {
		return c.String(500, "Error updating seller account balance")
	}

	return c.JSON(200, Response.SystemResponse{200, "Product bought successfully. The product you bought is",
		&echo.Map{"data": product}})
}

func AddMoneyToAccount(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer is called whether the parent/surrounding function is finished or not
	defer cancel()

	claims := configs.GetClaims(c)

	var user Model.User
	userID := claims.Id
	userIDObject, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.M{"_id": userIDObject}
	err := userCollection.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		return c.JSON(500, Response.SystemResponse{500, "User not found", &echo.Map{"data": err.Error()}})

	}

	var moneyString = c.Param("money")
	money, err := strconv.ParseFloat(moneyString, 64)
	user.AccountBalance = user.AccountBalance + money

	result, err1 := userCollection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})

	if err1 != nil || result.MatchedCount != 1 {
		return c.String(500, "Error updating user account balance")
	}

	return c.JSON(200, Response.SystemResponse{200, "Money added to account successfully. Your new account balance is",
		&echo.Map{"data": user.AccountBalance}})
}

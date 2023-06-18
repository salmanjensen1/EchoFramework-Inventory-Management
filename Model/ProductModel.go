package Model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	ProductName string             `json:"productname" validate:"required"`
	Quantity    int                `json:"quantity"`
	Price       float64            `json:"price,omitempty"`
	SellerID    string             `json:"sellerid,omitempty"`
}

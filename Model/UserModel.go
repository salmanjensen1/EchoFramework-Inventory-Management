package Model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `json:"name" validate:"required"`
	Username string             `json:"username" validate:"required"`
	Email    string             `json:"email" validate:"required"`
	Phone    string             `json:"phone" validate:"required"`
	Password string             `json:"password" validate:"required"`
	Address  string             `json:"address,omitempty"`
	IsAdmin  bool               `json:"isadmin,omitempty"`
}

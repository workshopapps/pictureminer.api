package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id, omitempty"`
	Username  string             `bson:"username" json:"username" validate:"required"`
	FirstName string             `bson:"first_name" json:"first_name"`
	LastName  string             `bson:"last_name" json:"last_name"`
	Email     string             `bson:"email" json:"email" validate:"required"`
	Password  string             `bson:"password" json:"password" validate:"required"`
}

type UserSignUpResponse struct {
	Username     string
	FirstName    string
	LastName     string
	Email        string
	Token        string
	TokenType    string
	ApiCallCount int
}

type UserLoginField struct {
	Email    string `bson:"email" json:"email" validate:"required"`
	Password string `bson:"password" json:"password" validate:"required"`
}

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

type UserStruct struct {
	Token        *string            `bson:"access_token" json:"access_token"`
	TokenType    string             `bson:"token_type" json:"token_type"`
	ID           primitive.ObjectID `bson:"_id"`
	UserName     string             `bson:"userName" json:"userName" validate:"required"`
	FirstName    string             `bson:"firstname" json:"firstname"`
	LastName     string             `bson:"lastname" json:"lastname"`
	Email        string             `bson:"email" json:"email" validate:"required"`
	Password     string             `bson:"password" json:"password" validate:"required"`
	ApiCallCount int                `bson:"api_call_count" json:"api_call_count"`
	//DataCreated  time.Time          `bson:"dataCreated" json:"datacreated"`
	//DateModified time.Time          `bson:"dateModified"json:"datemodified"`
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

package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

type UserLoginField struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

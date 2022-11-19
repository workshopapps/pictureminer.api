package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

//type User struct {
//	ID         primitive.ObjectID `bson:"_id"`
//	First_name *string            `bson:"first_name" validate:"required, min=5, max=100"`
//	Last_name  *string            `bson:"last_name" validate:"required, min=5, max=150"`
//	Email      *string            `bson:"email" validate:"email, required"`
//	Username   *string            `bson:"username" validate:"min=8"`
//	Password   *string            `bson:"password" validate:"required, min=8"`
//}

type UserStruct struct {
	UserName     primitive.ObjectID `bson:"userName" json:"userName" validate:"required, min=4, max=100"`
	Email        string             `bson:"email" json:"email" validate:"email, required"`
	Password     string             `bson:"password" json:"password" validate:"required, min 8"`
	Token        *string            `bson:"token" json:"token"`
	Request      []Requests         `bson:"request" json:"request"`
	DataCreated  time.Time          `bson:"dataCreated" json:"dataCreated"`
	DateModified time.Time          `bson:"dateModified"json:"dateModified"`
}
type Requests struct {
	ImageEncode  string    `json:"imageEncode" bson:"imageEncode"`
	DateCreated  time.Time `json:"dateCreated" bson:"dateCreated"`
	RequestsText string    `json:"requestsText" bson:"requestsText"`
}

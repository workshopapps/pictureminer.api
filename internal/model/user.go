package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id, omitempty"`
	Username    string             `bson:"username" json:"username" validate:"required"`
	FirstName   string             `bson:"first_name" json:"first_name"`
	LastName    string             `bson:"last_name" json:"last_name"`
	Email       string             `bson:"email" json:"email" validate:"required,email"`
	Password    string             `bson:"password" json:"password" validate:"required"`
	ProfileKey  string             `bson:"profile_key" json:"profile_key"`
	ProfileUrl  string             `bson:"profile_url" json:"profile_url"`
	DateCreated time.Time          `bson:"date_created" json:"date_created"`
	DateUpdated time.Time          `bson:"date_updated" json:" date_updated"`
}

type UserResponse struct {
	Username     string
	FirstName    string
	LastName     string
	Email        string
	ProfileKey   string
	ProfileUrl   string
	Token        string
	TokenType    string
	ApiCallCount int64
}

type UserLogin struct {
	Email    string `bson:"email" json:"email" validate:"required"`
	Password string `bson:"password" json:"password" validate:"required"`
}

// created this struct for swagger docs
type UserSignUp struct {
	Username    string             `bson:"username" json:"username" validate:"required"`
	FirstName   string             `bson:"first_name" json:"first_name"`
	LastName    string             `bson:"last_name" json:"last_name"`
	Email       string             `bson:"email" json:"email" validate:"required,email"`
	Password    string             `bson:"password" json:"password" validate:"required"`
}



type PasswordReset struct {
	Email           string `bson:"email" json:"email" validate:"required"`
	Password        string `bson:"password" json:"password" validate:"required"`
	ConfirmPassword string `bson:"confirm_password" json:"confirm_password" validate:"required"`
}

type PasswordForgot struct {
	Email string `bson:"email" json:"email" validate:"required"`
}

type UpdateUser struct {
	FirstName       string `bson:"first_name" json:"first_name"`
	LastName        string `bson:"last_name" json:"last_name"`
	Email           string `bson:"email" json:"email" validate:"required"`
	UserName        string `bson:"username" json:"username"`
	CurrentPassword string `bson:"current_password" json:"current_password"`
	NewPassword     string `bson:"new_password" json:"new_password"`
	ConfirmPassword string `bson:"confirm_password" json:"confirm_password"`
}

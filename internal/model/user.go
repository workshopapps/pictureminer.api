package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id, omitempty"`
	Username  string             `bson:"username" json:"username" validate:"required"`
	FirstName string             `bson:"first_name" json:"first_name"`
	LastName  string             `bson:"last_name" json:"last_name"`
	Email     string             `bson:"email" json:"email" validate:"required,email"`
	Plan      string             `bson:"plan" json:"plan"`
	//plan can be "free","starter" or "premium"
	Password     string    `bson:"password" json:"password" validate:"required"`
	ProfileKey   string    `bson:"profile_key" json:"profile_key"`
	ProfileUrl   string    `bson:"profile_url" json:"profile_url"`
	DateCreated  time.Time `bson:"date_created" json:"date_created"`
	DateUpdated  time.Time `bson:"date_updated" json:"date_updated"`
	ApiCallCount int64     `bson:"api_call_count" json:"api_call_count"`
	LastLogin    time.Time `bson:"last_login" json:"last_login"`
	IsVerified   bool      `bson:"is_verified" json:"is_verified"`
}

type UserResponse struct {
	Username     string
	FirstName    string
	LastName     string
	Email        string
	Plan         string
	ProfileKey   string
	ProfileUrl   string
	Token        string
	TokenType    string
	ApiCallCount int64
	LastLogin    time.Time
}

type UserLogin struct {
	Email    string `bson:"email" json:"email" validate:"required"`
	Password string `bson:"password" json:"password" validate:"required"`
}

// created this struct for swagger docs
type UserSignUp struct {
	Username  string `bson:"username" json:"username" validate:"required"`
	FirstName string `bson:"first_name" json:"first_name"`
	LastName  string `bson:"last_name" json:"last_name"`
	Email     string `bson:"email" json:"email" validate:"required,email"`
	Password  string `bson:"password" json:"password" validate:"required"`
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
	Email           string `bson:"email" json:"email"`
	UserName        string `bson:"username" json:"username"`
	CurrentPassword string `bson:"current_password" json:"current_password"`
	NewPassword     string `bson:"new_password" json:"new_password"`
	ConfirmPassword string `bson:"confirm_password" json:"confirm_password"`
}

type SubscriberEmail struct {
	ID               primitive.ObjectID `bson:"_id" json:"_id"`
	Email            string             `bson:"email" json:"email" validate:"required,email"`
	Subscribed       bool               `bson:"subscribed" json:"subscribed"`
	Price            float64            `bson:"price" json:"price"`
	SubscriptionType string             `bson:"subscription_type" json:"subscription_type"`
	ExpiresAt        time.Time          `bson:"expires_at" json:"expires_at"`
}

type SubscriptionRequest struct {
	Email            string  `bson:"email" json:"email" validate:"required"`
	Price            float64 `bson:"price" json:"price"`
	SubscriptionType string  `bson:"subscription_type" json:"subscription_type"`
}

package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MinedImage struct {
	ID           primitive.ObjectID `bson:"_id"`
	UserID       string             `bson:"user_id"`
	ImageKey     string             `bson:"image_key"`
	ImagePath    string             `bson:"image_path"`
	TextContent  string             `bson:"text_content"`
	DateCreated  time.Time          `bson:"date_created"`
	DateModified time.Time          `bson:"date_modified"`
}

type MinedImageResponse struct {
	ImagePath    string    `json:"image_path"`
	TextContent  string    `json:"text_content"`
	DateCreated  time.Time `json:"date_created"`
	DateModified time.Time `json:"date_modified"`
}

package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MinedImage struct {
	ID           primitive.ObjectID `bson:"_id"`
	UserID       string             `bson:"user_id"`
	ImageKey     string             `bson:"image_key"`
	ImageName    string             `bson:"image_name"`
	ImagePath    string             `bson:"image_path"`
	TextContent  string             `bson:"text_content"`
	DateCreated  time.Time          `bson:"date_created"`
	DateModified time.Time          `bson:"date_modified"`
}

type MineImageResponse struct {
	ImageName    string    `bson:"image_name" json:"image_name,omitempty"`
	ImagePath    string    `bson:"image_path" json:"image_path,omitempty"`
	TextContent  string    `bson:"text_content" json:"text_content"`
	DateCreated  time.Time `bson:"date_created" json:"date_created"`
	DateModified time.Time `bson:"date_modified" json:"date_modified"`
}

type MineImageUrlRequest struct {
	Url string `bson:"url" json:"url" validate:"url,required"`
}

type MineImagePromptRequest struct {
	Prompt string `form:"prompt" validate:"required"`
}

type MineImagePromptResponse struct {
	ImageName    string    `bson:"image_name" json:"image_name"`
	ImagePath    string    `bson:"image_path" json:"image_path"`
	TextPrompt   string    `bson:"text_prompt" json:"text_prompt"`
	CheckResult  bool      `bson:"check_result" json:"check_result"`
	DateCreated  time.Time `bson:"date_created" json:"date_created"`
	DateModified time.Time `bson:"date_modified" json:"date_modified"`
}


type BatchCollection struct {
	ID   primitive.ObjectID `bson:"batch_id" json:"batch_id"`
	User_id string `bson:"user_id" json:"user_id"`
	Name string `bson:"name" json:"name"`
	Description string `bson:"description" json:"tag"`
	Tags []string `bson:"tags" json:"tags"`
	Status string `bson:"status" json:"status"`
  }

type ImageCollection struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	Batch_id   string `bson:"batch_id" json:"batch_id"`
	Url string `bson:"url" json:"url"`
	Tag string `bson:"tag" json:"tag"`
  }

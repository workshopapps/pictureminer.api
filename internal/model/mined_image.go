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
	ImageKey     string    `bson:"image_key" json:"image_key,omitempty"`
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

type ProcessCallCount struct {
	MinedThisMonth    	int     `bson:"mined_this_month" json:"mined_this_month"`
	RemainingTomine   	int     `bson:"remaining_to_mine" json:"remaining_to_mine"`
	ImageCount    		int64   `bson:"image_count" json:"image_count"`
	BatchCount    		int64   `bson:"batch_count" json:"batch_count"`
	Status				bool	`bson:"status" json:"status"`
}

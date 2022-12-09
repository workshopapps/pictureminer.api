package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Batch struct {
	ID          primitive.ObjectID `bson:"_id"`
	UserID      string             `bson:"user_id"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Tags        []string           `bson:"tags"`
	Status      string             `bson:"status"`
	DateCreated time.Time          `bson:"date_created"`
}

type BatchResponse struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Tags        []string           `bson:"tags" json:"tags"`
	Status      string             `bson:"status" json:"status"`
	DateCreated time.Time          `bson:"date_created" json:"date_created"`
}

type BatchImage struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	BatchID string             `bson:"batch_id" json:"batch_id"`
	URL     string             `bson:"url" json:"url"`
	Tag     string             `bson:"tag" json:"tag"`
}

type ProcessBatchAPIRequest struct {
	Name        string   `json:"name,required"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,required"`
	Images      []string `json:"images"`
}

type ProcessBatchAPIResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name,required"`
	Description string    `json:"description,omitempty"`
	Tags        []string  `json:"tags,required"`
	Status      string    `json:"status"`
	DateCreated time.Time `bson:"date_created"`
}

type BatchDetailsResponse struct {
	Total    int `json:"total"`
	Tagged   int `json:"tagged"`
	Untagged int `json:"untagged"`
}

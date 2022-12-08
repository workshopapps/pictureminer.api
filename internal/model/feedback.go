package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)
type Feedback struct{
	ID     primitive.ObjectID `bson:"_id, omitempty"`
	ReviewerEmail string `bson:"reviewer_email" json:"reviewer_email"`
	ImageKey string `bson:"image_id" json:"image_key"`
	IsHelpful  bool `bson:"is_helpful" json:"is_helpful"`
	Feedback string `bson:"feedback" json:"feedback"`
	DateCreated  time.Time `bson:"date_created" json:"date_created"`
}

type FeedbackCreatedResponse struct {
	Message string `json:"message"`
}
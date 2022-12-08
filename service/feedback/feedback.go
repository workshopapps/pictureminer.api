package feedback

import (
	"context"
	"time"

	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/constants"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CollectFeedbackFromUser(feedback model.Feedback) (model.FeedbackCreatedResponse, error){
	
	feedback.ID = primitive.NewObjectID()
	feedback.DateCreated = time.Now()

	//save to the database
	database := config.GetConfig().Mongodb.Database
	feedbackCollection := mongodb.GetCollection(mongodb.Connection(), database, constants.FeedbackCollection)
	_, err := feedbackCollection.InsertOne(context.TODO(), feedback)

	if err != nil {
		return model.FeedbackCreatedResponse{} , err
	}
	feedbck:= model.FeedbackCreatedResponse{
		
		Message: "Response was successfully received",
	}
	return feedbck, nil
}

func GetAllFeedback()([]model.Feedback, error){
	ctx := context.TODO()
	cursor, err := mongodb.SelectFromCollection(ctx, config.GetConfig().Mongodb.Database, constants.FeedbackCollection, bson.M{})

	if err != nil {
		return []model.Feedback{}, err
	}

	var reviews = make([]model.Feedback, 0)
	cursor.All(ctx, &reviews)
	return reviews, nil
}
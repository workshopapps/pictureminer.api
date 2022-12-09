package user

import (
	"context"
	"time"

	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/constants"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SubscriberEmailResponse(subscriberEmail model.SubscriberEmail) (model.SubscriberEmail, string, int, error) {
	//create subscribtion time
	now:= time.Now()
	oneMonthLater:= now.AddDate(0,1,1)
	oneYearLater:= now.AddDate(1,0,0)
	
	if subscriberEmail.SubscriptionType == "MONTHLY" || subscriberEmail.SubscriptionType == "monthly"{
		subscriberEmail.ExpiresAt = oneMonthLater
	} else {
		subscriberEmail.ExpiresAt = oneYearLater
	}
	subscriberEmail.ID = primitive.NewObjectID()
	subscriberEmail.Email = subscriberEmail.Email
	subscriberEmail.Subscribed = true
	
	// save to DB
	database := config.GetConfig().Mongodb.Database
	subscriberCollection := mongodb.GetCollection(mongodb.Connection(), database, constants.SubscriberEmail)
	_, err := subscriberCollection.InsertOne(context.TODO(), subscriberEmail)
	if err != nil {
		return model.SubscriberEmail{}, "Unable to save email to database", 500, err
	}


	// build subscriber response
	subscriberResponse := model.SubscriberEmail{
		ID:     subscriberEmail.ID,
		Email:        subscriberEmail.Email,
		Subscribed:    true,
		ExpiresAt: subscriberEmail.ExpiresAt,
		SubscriptionType: subscriberEmail.SubscriptionType,
	}

	return subscriberResponse, "", 0, nil
}

func GetUserSubscribtion(email string)(model.SubscriberEmail, error){

	var subcribtion model.SubscriberEmail 
	//fetch from database
	database := config.GetConfig().Mongodb.Database
	ctx := context.TODO()
	subscriberCollection := mongodb.GetCollection(mongodb.Connection(), database, constants.SubscriberEmail)
	err :=subscriberCollection.FindOne(ctx, bson.M{"email": email}).Decode(&subcribtion)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return model.SubscriberEmail{}, err
		}
		return model.SubscriberEmail{}, err
	}
	return subcribtion, nil
}

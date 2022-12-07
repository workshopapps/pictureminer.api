package user

import (
	"context"
	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/constants"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SubscriberEmailResponse(subscriberEmail model.SubscriberEmail) (model.SubscriberEmail, string, int, error) {


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
	}

	return subscriberResponse, "", 0, nil
}

package user

import (
	"context"
	"errors"

	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/constants"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	batchservice "github.com/workshopapps/pictureminer.api/service/batch-service"
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
		ID:         subscriberEmail.ID,
		Email:      subscriberEmail.Email,
		Subscribed: true,
	}

	return subscriberResponse, "", 0, nil
}

func VerifyEmail(userID interface{}) error {
	id, ok := userID.(string)
	if !ok {
		return errors.New("Invalid user ID")
	}

	response, err := mongodb.MongoUpdate(id[10:len(id)-2], map[string]interface{}{
		"is_verified": true,
	}, constants.UserCollection)
	if err != nil {
		return err
	}

	if response.MatchedCount == 0 {
		return errors.New("User ID not found")
	}

	return nil
}

func IsVerified(userID string) (bool, error) {
	user, _, err := batchservice.GetUserFromID(userID)
	if err != nil {
		return false, err
	}

	return user.IsVerified, nil
}

package user

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/constants"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/s3"
	"github.com/workshopapps/pictureminer.api/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func SubscriberEmailResponse(subscriberEmail model.SubscriberEmail) (model.SubscriberEmail, string, int, error) {


	subscriberEmail.ID = primitive.NewObjectID()
	subscriberEmail.Email = email
	subscriberEmail.Subscribed = true

	// save to DB
	database := config.GetConfig().Mongodb.Database
	subscriberCollection := mongodb.GetCollection(mongodb.Connection(), database, constants.SubscriberEmail)
	_, err = subscriberCollection.InsertOne(context.TODO(), subscriberEmail)
	if err != nil {
		return model.UserResponse{}, "Unable to save email to database", 500, err
	}

	secretkey := config.GetConfig().Server.Secret
	token, err := utility.CreateToken("id", user.ID.String(), secretkey)
	if err != nil {
		return model.UserResponse{}, fmt.Sprintf("unable to create token: %v", err.Error()), 500, err
	}

	// build user response
	userResponse := model.EmailResponse{
		ID:     user.Username,
		Email:        user.Email,
		Subscribed:        user.Email
	}

	return userResponse, "", 0, nil
}


https://dictionary.cambridge.org/images/thumb/cushio_noun_001_04014.jpg

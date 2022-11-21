package user

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/go-playground/validator/v10"
	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/constants"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"github.com/workshopapps/pictureminer.api/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func SignUpUser(user model.User) (model.UserResponse, string, int, error) {
	// check if user already exists
	_, err := getUserFromDB(user.Email)
	if err == nil {
		return model.UserResponse{}, "user already exist", 403, validator.ValidationErrors{}
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password = string(hash)
	user.ID = primitive.NewObjectID()

	// save to DB
	database := config.GetConfig().Mongodb.Database
	userCollection := mongodb.GetCollection(mongodb.Connection(), database, constants.UserCollection)
	_, err = userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return model.UserResponse{}, "Unable to save user to database", 500, validator.ValidationErrors{}
	}

	secretkey := config.GetConfig().Server.Secret
	token, err := utility.CreateToken("id", user.ID.String(), secretkey)
	if err != nil {
		return model.UserResponse{}, fmt.Sprintf("unable to create token: %v", err.Error()), 500, validator.ValidationErrors{}
	}

	// build user response
	userResponse := model.UserResponse{
		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		TokenType:    "bearer",
		Token:        token,
		ApiCallCount: 0,
	}

	return userResponse, "", 0, nil
}

func LoginUser(userLoginObject model.UserLogin) (model.UserResponse, string, int, error) {
	user, err := getUserFromDB(userLoginObject.Email)
	if err != nil {
		return model.UserResponse{}, "user does not exist", 404, validator.ValidationErrors{}
	}

	if !isValidPassword(user.Password, userLoginObject.Password) {
		return model.UserResponse{}, "invalid password", 401, validator.ValidationErrors{}
	}

	secretkey := config.GetConfig().Server.Secret
	token, err := utility.CreateToken("id", user.ID.String(), secretkey)
	if err != nil {
		return model.UserResponse{}, fmt.Sprintf("unable to create token: %v", err.Error()), 500, validator.ValidationErrors{}
	}

	// build user response
	userResponse := model.UserResponse{
		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		TokenType:    "bearer",
		Token:        token,
		ApiCallCount: rand.Intn(10),
	}

	return userResponse, "", 0, nil
}

func getUserFromDB(email string) (model.User, error) {
	var user model.User
	database := config.GetConfig().Mongodb.Database
	userCollection := mongodb.GetCollection(mongodb.Connection(), database, constants.UserCollection)
	result := userCollection.FindOne(context.TODO(), bson.M{"email": email})
	err := result.Err()
	if err != nil {
		return model.User{}, err
	}

	err = result.Decode(&user)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func isValidPassword(userPassword, providedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword)) == nil
}

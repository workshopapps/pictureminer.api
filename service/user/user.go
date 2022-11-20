package user

import (
	"context"
	"math/rand"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/workshopapps/pictureminer.api/internal/constants"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"github.com/workshopapps/pictureminer.api/utility"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func PasswordIsValid(userPassword, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))

	check := true
	msg := ""

	if err != nil {
		check = false
		msg = "invalid email or password"
	}
	return check, msg
}

func CheckUserExists(user model.UserLoginField) (model.User, error) {
	var profile model.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userCollection := mongodb.GetCollection(mongodb.ConnectToDB(), constants.UserDatabase, constants.UserCollection)
	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&profile)
	if err != nil {
		return profile, err
	}
	return profile, nil
}

func SignUpUser(user model.User) (model.UserSignUpResponse, string, error) {
	//Database connection
	mongoClient := mongodb.Connection()
	userCollection := mongodb.GetCollection(mongoClient, constants.UserDatabase, constants.UserCollection)

	//	Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return model.UserSignUpResponse{}, "Unable to hash password", err
	}
	user.Password = string(hash)

	// save to DB
	userCollection.InsertOne(context.Background(), user)

	// build user response
	var userResponse model.UserSignUpResponse
	userResponse.Username = user.Username
	userResponse.FirstName = user.FirstName
	userResponse.LastName = user.LastName
	userResponse.Email = user.Email

	userResponse.TokenType = "bearer"
	token, err := utility.CreateToken(userResponse.Email)
	if err != nil {
		return model.UserSignUpResponse{}, "unable to create access token", err

	}
	userResponse.Token = token
	userResponse.ApiCallCount = 0

	return userResponse, "", nil
}

func LoginUser(userLoginObject model.UserLoginField) (model.UserSignUpResponse, string, error) {
	//Database connection
	mongoClient := mongodb.Connection()
	userCollection := mongodb.GetCollection(mongoClient, constants.UserDatabase, constants.UserCollection)

	var user model.User
	ctx := context.TODO()
	result := userCollection.FindOne(ctx, bson.M{"email": userLoginObject.Email})
	err := result.Err()
	if err != nil {
		return model.UserSignUpResponse{}, "unable to find user", validator.ValidationErrors{}
	}

	err = result.Decode(&user)
	if err != nil {
		return model.UserSignUpResponse{}, "unable to decode user", validator.ValidationErrors{}
	}

	// build user response
	var userResponse model.UserSignUpResponse
	userResponse.Username = user.Username
	userResponse.FirstName = user.FirstName
	userResponse.LastName = user.LastName
	userResponse.Email = user.Email

	userResponse.TokenType = "bearer"
	token, err := utility.CreateToken(userResponse.Email)
	if err != nil {
		return model.UserSignUpResponse{}, "unable to create access token", validator.ValidationErrors{}

	}
	userResponse.Token = token
	userResponse.ApiCallCount = rand.Intn(10)

	return userResponse, "", nil
}

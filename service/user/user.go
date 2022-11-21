package user

import (
	"context"
	"math/rand"

	"github.com/go-playground/validator/v10"
	"github.com/workshopapps/pictureminer.api/internal/constants"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"github.com/workshopapps/pictureminer.api/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func GetUserFromDB(email string) (model.User, error) {
	var user model.User
	userCollection := mongodb.GetCollection(mongodb.ConnectToDB(), constants.UserDatabase, constants.UserCollection)
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

func SignUpUser(user model.User) (model.UserSignUpResponse, string, int, error) {
	_, err := GetUserFromDB(user.Email)
	if err == nil {
		return model.UserSignUpResponse{}, "user already exist", 403, validator.ValidationErrors{}
	}

	//Database connection
	userCollection := mongodb.GetCollection(mongodb.ConnectToDB(), constants.UserDatabase, constants.UserCollection)

	//	Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return model.UserSignUpResponse{}, "Unable to hash password", 500, err
	}
	user.Password = string(hash)
	user.ID = primitive.NewObjectID()

	// save to DB
	_, err = userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return model.UserSignUpResponse{}, "Unable to save user to database", 500, validator.ValidationErrors{}
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
		return model.UserSignUpResponse{}, "unable to create access token", 500, err

	}
	userResponse.Token = token
	userResponse.ApiCallCount = 0

	return userResponse, "", 0, nil
}

func isValidPassword(userPassword, providedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	if err != nil {
		return false
	}
	return true
}

func LoginUser(userLoginObject model.UserLoginField) (model.UserSignUpResponse, string, int, error) {
	user, err := GetUserFromDB(userLoginObject.Email)
	if err != nil {
		return model.UserSignUpResponse{}, "user does not exist", 404, validator.ValidationErrors{}
	}

	if !isValidPassword(user.Password, userLoginObject.Password) {
		return model.UserSignUpResponse{}, "Invalid password", 401, validator.ValidationErrors{}
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
		return model.UserSignUpResponse{}, err.Error(), 500, validator.ValidationErrors{}

	}
	userResponse.Token = token
	userResponse.ApiCallCount = rand.Intn(10)

	return userResponse, "", 0, nil
}

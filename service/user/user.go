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

func SignUpUser(user model.User) (model.UserSignUpResponse, string, error) {
	_, err := GetUserFromDB(user.Email)
	if err == nil {
		return model.UserSignUpResponse{}, "user already exist", validator.ValidationErrors{}
	}

	//Database connection
	userCollection := mongodb.GetCollection(mongodb.ConnectToDB(), constants.UserDatabase, constants.UserCollection)

	//	Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return model.UserSignUpResponse{}, "Unable to hash password", err
	}
	user.Password = string(hash)
	user.ID = primitive.NewObjectID()

	// save to DB
	_, err = userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return model.UserSignUpResponse{}, "Unable to save user to database", validator.ValidationErrors{}
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
		return model.UserSignUpResponse{}, "unable to create access token", err

	}
	userResponse.Token = token
	userResponse.ApiCallCount = 0

	return userResponse, "", nil
}

func LoginUser(userLoginObject model.UserLoginField) (model.UserSignUpResponse, string, error) {
	user, err := GetUserFromDB(userLoginObject.Email)
	if err != nil {
		return model.UserSignUpResponse{}, err.Error(), validator.ValidationErrors{}
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
		return model.UserSignUpResponse{}, err.Error(), validator.ValidationErrors{}

	}
	userResponse.Token = token
	userResponse.ApiCallCount = rand.Intn(10)

	return userResponse, "", nil
}

package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/constants"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"github.com/workshopapps/pictureminer.api/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func SignUpUser(user model.User) (model.UserResponse, string, int, error) {
	// check if user already exists
	_, err := getUserFromDB(user.Email)
	if err == nil {
		return model.UserResponse{}, "user already exist", 403, errors.New("user already exist in database")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password = string(hash)
	user.ID = primitive.NewObjectID()

	// save to DB
	database := config.GetConfig().Mongodb.Database
	userCollection := mongodb.GetCollection(mongodb.Connection(), database, constants.UserCollection)
	_, err = userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return model.UserResponse{}, "Unable to save user to database", 500, err
	}

	secretkey := config.GetConfig().Server.Secret
	token, err := utility.CreateToken("id", user.ID.String(), secretkey)
	if err != nil {
		return model.UserResponse{}, fmt.Sprintf("unable to create token: %v", err.Error()), 500, err
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
		return model.UserResponse{}, "user does not exist", 404, err
	}

	if !isValidPassword(user.Password, userLoginObject.Password) {
		return model.UserResponse{}, "invalid password", 401, errors.New("invalid password")
	}

	secretkey := config.GetConfig().Server.Secret
	token, err := utility.CreateToken("id", user.ID.String(), secretkey)
	if err != nil {
		return model.UserResponse{}, fmt.Sprintf("unable to create token: %v", err.Error()), 500, err
	}

	// implementaton code
	estCount, err := mongodb.CountFromCollection(user.ID)
	if err != nil {
		return model.UserResponse{}, "error reading number of documents", 500, err
	}

	// build user response
	userResponse := model.UserResponse{
		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		TokenType:    "bearer",
		Token:        token,
		ApiCallCount: estCount,
	}

	return userResponse, "", 0, nil
}

func ResetPassword(reqBody model.PasswordReset) (int, error) {
	user, err := getUserFromDB(reqBody.Email)
	if err != nil {
		return 404, fmt.Errorf("user does not exist: %s", err.Error())
	}

	if reqBody.Password != reqBody.ConfirmPassword {
		return 401, errors.New("passwords do not match")
	}

	newPasswordHash, _ := bcrypt.GenerateFromPassword([]byte(reqBody.Password), 10)

	// update user in db
	database := config.GetConfig().Mongodb.Database
	userCollection := mongodb.GetCollection(mongodb.Connection(), database, constants.UserCollection)
	filter := bson.M{"email": user.Email}
	update := bson.D{{"$set", bson.D{{"password", newPasswordHash}}}}
	_, err = userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return 500, fmt.Errorf("unable to update user password: %s", err.Error())
	}

	return 0, nil
}

func ForgotPassword(reqBody model.PasswordForgot) (int, error) {
	_, err := getUserFromDB(reqBody.Email)
	if err != nil {
		return 404, fmt.Errorf("user does not exist: %s", err.Error())
	}
	
	var w http.ResponseWriter
	var r *http.Request

	http.Redirect(w, r, "/reset", http.StatusFound)
	return http.StatusOK, nil
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


func UpdateUserService(user model.UpdateUser) (*mongo.UpdateResult, error) {
	database := config.GetConfig().Mongodb.Database
	userCollection := mongodb.GetCollection(mongodb.Connection(), database, constants.UserCollection)

	filter := bson.M{"email": user.Email}

	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "first_name", Value: user.FirstName},
		{Key: "last_name", Value: user.LastName},
		{Key: "email", Value: user.Email},
		{Key: "username", Value: user.UserName}} }}

	result, err := userCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return nil, err
	}

	return result,nil
}
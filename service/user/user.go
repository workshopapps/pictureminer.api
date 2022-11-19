package user

import (
	"context"
	"time"

	"github.com/workshopapps/pictureminer.api/internal/constants"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
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

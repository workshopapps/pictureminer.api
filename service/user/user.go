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

var (
	FreePlan    = "free"
	StarterPlan = "starter"
	PremiumPlan = "premium"
)

func SignUpUser(user model.User) (model.UserResponse, string, int, error) {
	// check if user already exists
	_, err := getUserFromDB(user.Email)
	if err == nil {
		return model.UserResponse{}, "user already exist", 403, errors.New("user already exist in database")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	profile_url, profile_key := s3.DefaultProfile()

	user.Password = string(hash)
	user.ID = primitive.NewObjectID()
	user.ProfileUrl = profile_url
	user.ProfileKey = profile_key
	user.Plan = FreePlan

	// save to DB
	_, err = mongodb.MongoPost(constants.UserCollection, user)
	if err != nil {
		return model.UserResponse{}, "Unable to save user to database", 500, err
	}

	secretkey := config.GetConfig().Server.Secret
	token, err := utility.CreateToken("id", user.ID.String(), secretkey)
	if err != nil {
		return model.UserResponse{}, fmt.Sprintf("unable to create token: %v", err.Error()), 500, err
	}

	// Send verification email
	sendVerficationEmail(token, &user)

	// build user response
	userResponse := model.UserResponse{
		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		Plan:         user.Plan,
		ProfileKey:   user.ProfileKey,
		ProfileUrl:   user.ProfileUrl,
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

	//update user's last login time
	database := config.GetConfig().Mongodb.Database
	userCollection := mongodb.GetCollection(mongodb.Connection(), database, constants.UserCollection)
	filter := bson.M{"email": user.Email}
	update := bson.D{{"$set", bson.D{{"last_login", time.Now()}}}}
	_, err = userCollection.UpdateOne(context.TODO(), filter, update)

	// build user response
	userResponse := model.UserResponse{
		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		Plan:         user.Plan,
		ProfileKey:   user.ProfileKey,
		ProfileUrl:   user.ProfileUrl,
		TokenType:    "bearer",
		Token:        token,
		ApiCallCount: user.ApiCallCount,
		LastLogin:    user.LastLogin,
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

func ProfilePictureServiceUpload(userId interface{}, image io.ReadCloser, filename string) (string, error) {
	string_id, ok := userId.(string)
	if !ok {
		return "", errors.New("invalid userid")
	}

	id := string_id[10 : len(string_id)-2]
	imagePath, err := s3.UploadProfileImage(image, id+filepath.Ext(filename))
	if err != nil {
		return "", err
	}

	updateUserPicture := map[string]interface{}{
		"profile_url":  imagePath,
		"date_updated": time.Now(),
	}

	update_response, err := mongodb.MongoUpdate(id, updateUserPicture, constants.UserCollection)
	if err != nil {
		return "", err
	}

	if update_response.MatchedCount != 1 {
		return "", fmt.Errorf("User with ID not found")
	}

	return imagePath, nil

}

func UpdateUserService(user model.UpdateUser, userId interface{}) (int, error) {
	string_id, _ := userId.(string)
	string_id = string_id[10 : len(string_id)-2]
	id, _ := primitive.ObjectIDFromHex(string_id)
	database := config.GetConfig().Mongodb.Database
	filter := bson.D{{"_id", id}}

	if len(user.LastName) > 0 {
		update := bson.D{
			{"$set", bson.D{{"last_name", user.LastName}}},
		}
		err := UpdateFunc(database, filter, update)
		if err != nil {
			return 400, err
		}
	}

	if len(user.FirstName) > 0 {
		update := bson.D{
			{"$set", bson.D{{"first_name", user.FirstName}}},
		}
		err := UpdateFunc(database, filter, update)
		if err != nil {
			return 400, err
		}
	}

	if len(user.Email) > 0 {
		update := bson.D{
			{"$set", bson.D{{"email", user.Email}}},
		}
		err := UpdateFunc(database, filter, update)
		if err != nil {
			return 400, err
		}
	}

	if len(user.UserName) > 0 {
		update := bson.D{
			{"$set", bson.D{{"username", user.UserName}}},
		}
		err := UpdateFunc(database, filter, update)
		if err != nil {
			return 400, err
		}
	}

	if len(user.NewPassword) > 0 {
		err := CheckPasswords(user, id)
		if err != nil {
			return 400, err
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(user.NewPassword), 10)
		user.NewPassword = string(hash)
		update := bson.D{
			{"$set", bson.D{{"password", user.NewPassword}}},
		}
		err = UpdateFunc(database, filter, update)
		if err != nil {
			return 400, err
		}
	}

	return 200, nil
}

func UpdateFunc(db string, filter bson.D, update bson.D) error {
	userCollection := mongodb.GetCollection(mongodb.Connection(), db, constants.UserCollection)
	_, err := userCollection.UpdateOne(context.TODO(), filter, update)
	return err
}

func CheckPasswords(user model.UpdateUser, id primitive.ObjectID) error {
	var userDocument model.User
	database := config.GetConfig().Mongodb.Database
	userCollection := mongodb.GetCollection(mongodb.Connection(), database, constants.UserCollection)
	result := userCollection.FindOne(context.TODO(), bson.M{"_id": id})

	err := result.Decode(&userDocument)
	if err != nil {
		return err
	}
	if len(user.CurrentPassword) < 0 {
		err := errors.New("Provide current password")
		return err
	}
	if !isValidPassword(userDocument.Password, user.CurrentPassword) {
		err := errors.New("Password invalid")
		return err
	}
	if user.NewPassword != user.ConfirmPassword {
		err := errors.New("Passwords do not match")
		return err
	}
	return err
}

func sendVerficationEmail(token string, user *model.User) error {
	from := config.GetConfig().NotifyEmail.Email
	password := config.GetConfig().NotifyEmail.Password
	username := config.GetConfig().NotifyEmail.Username

	data := &utility.EmailData{
		Subject:  "Verify Email for Discripto",
		URL:      constants.VerifyEmailURL + "?token=" + token,
		UserName: user.Username,
	}

	if err := utility.SendMailNaked(from, username, password, user.Email, constants.VerifyEmailTemplate, data); err != nil {
		return err
	}

	return nil
}

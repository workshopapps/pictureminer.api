package user

import (
	"context"
	"fmt"
	"time"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/workshopapps/pictureminer.api/internal/model"

	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"github.com/workshopapps/pictureminer.api/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func (base *Controller) CreateUser(c *gin.Context) {

	rd := utility.BuildSuccessResponse(http.StatusCreated, "user created successfully", gin.H{"user": "user object"})
	c.JSON(http.StatusOK, rd)

}

// getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("ImageCollection").Collection(collectionName)
	return collection
}

func passwordIsValid(userPassword, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))

	check := true
	msg := ""

	if err != nil {
		check = false
		msg = "invalid email or password"
	}
	return check, msg
}

func (base *Controller) Login(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userCollection := GetCollection(mongodb.ConnectToDB(), "users")

	var user UserLoginField
	var profile model.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error a": err.Error()})
		return
	}

	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&profile)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "404 Not Found", "User does not exits", err.Error(), fmt.Sprintf("%s does not exist", user.Email))
		c.JSON(http.StatusUnauthorized, rd)
		return

	}

	//isValid, msg := passwordIsValid(*profile.Password, user.Password)

	//if isValid != true {
	//c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
	//return
	//} else {
	//rd := utility.BuildSuccessResponse(http.StatusOK, "user logged successfully", gin.H{"user": user.Email})
	//c.JSON(http.StatusOK, rd)

	//}

	if user.Email == *profile.Email && user.Password == *profile.Password {
		rd := utility.BuildSuccessResponse(http.StatusOK, "user logged successfully", gin.H{"user": user.Email})
		c.JSON(http.StatusOK, rd)

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
	}
	rd := utility.BuildSuccessResponse(http.StatusCreated, "user created successfully", gin.H{"user": "login object"})
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) Signup(c *gin.Context) {

	//Database connection
	mongoClient := mongodb.Connection()
	imageDB := mongoClient.Database("ImageCollection")
	userCollection := imageDB.Collection("user")

	//Binding the userdetails to userStruct
	var User model.UserStruct
	err := c.Bind(&User)
	if err != nil {
		log.Fatal("Unable to bind user signup details")
	}
	//	Hashing the password
	harsh, err := bcrypt.GenerateFromPassword([]byte(User.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to harsh paswword",
		})
		return
	}
	User.Password = string(harsh)

	userCollection.InsertOne(context.Background(), User)
	c.JSON(200, User)

}

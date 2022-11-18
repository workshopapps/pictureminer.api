package user

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/workshopapps/pictureminer.api/internal/model"
	DB "github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongo"
	"github.com/workshopapps/pictureminer.api/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// "github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongo"
func (base *Controller) CreateUser(c *gin.Context) {

	rd := utility.BuildSuccessResponse(http.StatusCreated, "user created successfully", gin.H{"user": "user object"})
	c.JSON(http.StatusOK, rd)

}

// getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("pixminer").Collection(collectionName)
	return collection
}

//func isPasswordValid(userPassword, providedPassword string) (bool, string) {
//err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))

//check := true
//msg := ""

//if err != nil {
//check = false
//msg = "invalid username or password"
//}

//return check, msg
//}

func (base *Controller) Login(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userCollection := GetCollection(DB.ConnectToDB(), "users")

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

	if user.Email == *profile.Email && user.Password == *profile.Password {
		rd := utility.BuildSuccessResponse(http.StatusOK, "user logged successfully", gin.H{"user": user.Email})
		c.JSON(http.StatusOK, rd)

	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
	}

}

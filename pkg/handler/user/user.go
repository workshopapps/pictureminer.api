package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"github.com/workshopapps/pictureminer.api/utility"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func (base *Controller) CreateUser(c *gin.Context) {

	rd := utility.BuildSuccessResponse(http.StatusCreated, "user created successfully", gin.H{"user": "user object"})
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) Login(c *gin.Context) {

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
	User.Token = utility.CreateToken(&User)

	userCollection.InsertOne(context.Background(), User)
	c.JSON(200, User)

}

package user

import (
	"context"
	"fmt"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/workshopapps/pictureminer.api/internal/model"
	usService "github.com/workshopapps/pictureminer.api/service/user"

	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"github.com/workshopapps/pictureminer.api/utility"
	"golang.org/x/crypto/bcrypt"
)

func (base *Controller) CreateUser(c *gin.Context) {

	rd := utility.BuildSuccessResponse(http.StatusCreated, "user created successfully", gin.H{"user": "user object"})
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) Login(c *gin.Context) {
	var user model.UserLoginField

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := usService.CheckUserExists(user)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "404 Not Found", "User does not exits", err.Error(), fmt.Sprintf("%s does not exist", user.Email))
		c.JSON(http.StatusUnauthorized, rd)
		return

	}

	if user.Email == *profile.Email {

		isValid, msg := usService.PasswordIsValid(*profile.Password, user.Password)
		if isValid != true {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}

		rd := utility.BuildSuccessResponse(http.StatusOK, "user logged successfully", gin.H{"user": user.Email})
		c.JSON(http.StatusOK, rd)

	}

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
	User.Token = utility.CreateToken()

	userCollection.InsertOne(context.Background(), User)
	c.JSON(200, User)

}

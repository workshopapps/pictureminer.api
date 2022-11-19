package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	userService "github.com/workshopapps/pictureminer.api/service/user"
	"github.com/workshopapps/pictureminer.api/utility"
	"golang.org/x/crypto/bcrypt"
)

func (base *Controller) CreateUser(c *gin.Context) {

	rd := utility.BuildSuccessResponse(http.StatusCreated, "user created successfully", gin.H{"user": "user object"})
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
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Unable to bind user signup details", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validate.Struct(&User)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validate), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	//	Hashing the password
	harsh, err := bcrypt.GenerateFromPassword([]byte(User.Password), 10)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to harsh paswword", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	User.Password = string(harsh)
	User.TokenType = "bearer"
	token, err := utility.CreateToken(User.Email)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return

	}
	User.Token = &token
	User.ApiCallCount = 0
	//User.ID = ObjectID()

	userCollection.InsertOne(context.Background(), User)
	object := utility.BuildSuccessResponse(200, "User created successfully", User)
	c.JSON(200, object)
}

func (base *Controller) Login(c *gin.Context) {
	var user model.UserLoginField

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := userService.CheckUserExists(user)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusNotFound, "404 Not Found", "User does not exits", err.Error(), fmt.Sprintf("%s does not exist", user.Email))
		c.JSON(http.StatusUnauthorized, rd)
		return

	}

	if user.Email == profile.Email {

		isValid, msg := userService.PasswordIsValid(profile.Password, user.Password)
		if isValid != true {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}

		rd := utility.BuildSuccessResponse(http.StatusOK, "user logged successfully", gin.H{"token_type": profile.TokenType, "username": profile.UserName, "email": profile.Email, "firstname": profile.FirstName, "lastname": profile.LastName, "api_call_count": profile.ApiCallCount})
		c.JSON(http.StatusOK, rd)

	}

}

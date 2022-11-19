package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/workshopapps/pictureminer.api/internal/constants"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"github.com/workshopapps/pictureminer.api/utility"
	"golang.org/x/crypto/bcrypt"
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
	userCollection := mongodb.GetCollection(mongoClient, constants.UserDatabase, constants.UserCollection)

	//Binding the userdetails to userStruct
	var User model.User
	err := c.Bind(&User)
	fmt.Println(User)
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
	hash, err := bcrypt.GenerateFromPassword([]byte(User.Password), 10)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Failed to hash paswword", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	User.Password = string(hash)
	userCollection.InsertOne(context.Background(), User)

	var userResponse model.UserSignUpResponse
	userResponse.Username = User.Username
	userResponse.FirstName = User.FirstName
	userResponse.LastName = User.LastName
	userResponse.Email = User.Email

	userResponse.TokenType = "bearer"
	token, err := utility.CreateToken(userResponse.Email)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return

	}
	userResponse.Token = token
	userResponse.ApiCallCount = 0

	object := utility.BuildSuccessResponse(200, "User created successfully", userResponse)
	c.JSON(200, object)
}

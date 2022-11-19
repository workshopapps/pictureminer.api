package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
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

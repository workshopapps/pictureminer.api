package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/service/user"
	"github.com/workshopapps/pictureminer.api/utility"
)

func (base *Controller) CreateUser(c *gin.Context) {

	rd := utility.BuildSuccessResponse(http.StatusCreated, "user created successfully", gin.H{"user": "user object"})
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) Signup(c *gin.Context) {

	// bind userdetails to User struct
	var User model.User
	err := c.Bind(&User)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Unable to bind user signup details", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validate.Struct(&User)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validate), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userResponse, msg, code, err := user.SignUpUser(User)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", msg, err, nil)
		c.JSON(code, rd)
		return
	}

	object := utility.BuildSuccessResponse(200, "signup successful", userResponse)
	c.JSON(200, object)
}

func (base *Controller) Login(c *gin.Context) {
	// bind user login details to User struct
	var User model.UserLogin
	err := c.Bind(&User)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Unable to bind user login details", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validate.Struct(&User)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validate), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	userResponse, msg, code, err := user.LoginUser(User)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", msg, err, nil)
		c.JSON(code, rd)
		return
	}

	object := utility.BuildSuccessResponse(200, "User login successful", userResponse)
	c.JSON(200, object)
}

func (base *Controller) ResetPassword(c *gin.Context) {
	// bind password reset details to User struct
	var reqBody model.PasswordReset
	err := c.Bind(&reqBody)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Unable to bind password reset details", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validate.Struct(&reqBody)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validate), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	code, err := user.ResetPassword(reqBody)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", "password reset failed", gin.H{"error": err.Error()}, nil)
		c.JSON(code, rd)
		return
	}

	object := utility.BuildSuccessResponse(200, "password reset success", gin.H{})
	c.JSON(200, object)
}

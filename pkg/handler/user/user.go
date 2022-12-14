package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/service/user"
	"github.com/workshopapps/pictureminer.api/utility"
)

func (base *Controller) CreateUser(c *gin.Context) {

	rd := utility.BuildSuccessResponse(http.StatusCreated, "user created successfully", gin.H{"user": "user object"})
	c.JSON(http.StatusOK, rd)

}

// Signup          godoc
// @Summary		Signs Up a User
// @Description Creates an account for a new user
// @Tags        users
// @Produce     json
// @Param User body model.UserSignUp true "User Signup" model.User
// @Success     200  {object} model.UserResponse
// @Router      /signup [post]
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

// Login          godoc
// @Summary		Login User
// @Description Logs in a User
// @Tags        users
// @Produce     json
// @Param User body model.UserLogin true "User Login" model.UserLogin
// @Success     200  {object} model.UserLogin
// @Router      /login [post]
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

// Post             godoc
// @Summary     Resests the password of the user
// @Description Send a post request to reset th password of the user
// @Tags        users
// @Produce     json
// @Param       ping body     model.PasswordReset true "Ping JSON"
// @Success     200  {object} utility.Response
// @Router      /reset [post]
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

// Post             godoc
// @Summary     Checks the status of the forgot passoword
// @Description Send a dummy post request to test the status of the server
// @Tags        Forgot Password
// @Produce     json
// @Param       ping body     model.PasswordForgot true "Ping JSON"
// @Success     200  {object} utility.Response
// @Router      /forgot-password [post]
func (base *Controller) ForgotPassword(c *gin.Context) {
	// validate jwt token
	// secretKey := config.GetConfig().Server.Secret
	// token := utility.ExtractToken(c)
	// _, err := utility.GetKey("id", token, secretKey)
	// if err != nil {
	// 	rd := utility.BuildErrorResponse(http.StatusUnauthorized, "failed", "could not verify token", gin.H{"error": err.Error()}, nil)
	// 	c.JSON(http.StatusUnauthorized, rd)
	// 	return
	// }
	// bind password reset details to User struct
	var reqBody model.PasswordForgot
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

	code, err := user.ForgotPassword(reqBody)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", "password reset failed", gin.H{"error": err.Error()}, nil)
		c.JSON(code, rd)
		return
	}

	object := utility.BuildSuccessResponse(200, "password reset success", gin.H{})
	c.JSON(200, object)
}

// Update Profile Pic             godoc
// @Summary     Updates a User profile picture image
// @Description Send a patch request containing a file to be updated and receives a response of its url path after upload.
// @Tags        users
// @Param       image formData file true "image"
// @Success     200  {object} utility.Response
// @Router      /update_user_picture [patch]
// @Security BearerAuth
func (base *Controller) UpdateProfilePicture(c *gin.Context) {

	secretKey := config.GetConfig().Server.Secret
	token := utility.ExtractToken(c)
	userId, err := utility.GetKey("id", token, secretKey)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, "failed", "could not verify token", nil, gin.H{"error": err.Error()})
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

	if c.ContentType() != "multipart/form-data" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "invalid request", nil, gin.H{"error": "file is not present"})
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	image, imageHeader, err := c.Request.FormFile("image")
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "could not parse file", nil, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	defer image.Close()

	if !utility.ValidImageFormat(imageHeader.Filename) {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "invalid file", nil, gin.H{"error": "file is not an image"})
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	picturePath, err := user.ProfilePictureServiceUpload(userId, image, imageHeader.Filename)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "undefined error", nil, err.Error())
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, "Profile picture upload successful", gin.H{"profile_url": picturePath})
	c.JSON(http.StatusOK, rd)
}

// Update User          godoc
// @Summary		Update User
// @Description Updates a User's information - email,firstName,lastName,password - Bearer token required - To change password, current_password, new_password and confirm_password(repeat of the new password) are required
// @Tags        users
// @Produce     json
// @Param User body model.UpdateUser true "User Update" model.UserUpdate
// @Success     200
// @Router      /update-user [patch]
func (base *Controller) UpdateUser(c *gin.Context) {
	secretKey := config.GetConfig().Server.Secret
	token := utility.ExtractToken(c)
	userId, err := utility.GetKey("id", token, secretKey)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, "failed", "could not verify token", nil, gin.H{"error": err.Error()})
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

	var reqBody model.UpdateUser
	if err = c.Bind(&reqBody); err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Unable to bind user update details", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validate.Struct(&reqBody)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validate), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	statusCode, err := user.UpdateUserService(reqBody, userId)
	if err != nil {
		rd := utility.BuildErrorResponse(statusCode, "error", "user update failed", gin.H{"error": err.Error()}, nil)
		c.JSON(statusCode, rd)
		return
	}

	object := utility.BuildSuccessResponse(200, "User update successful", gin.H{})
	c.JSON(statusCode, object)
}

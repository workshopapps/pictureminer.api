package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/service/admin"
	"github.com/workshopapps/pictureminer.api/utility"
)

type Controller struct {
	Validate *validator.Validate
	Logger   *utility.Logger
}

// Getusers
func (base *Controller) GetUsers(c *gin.Context) {

	users, err := admin.GetUsers()
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", err.Error(), err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "success", gin.H{"data": users})
	c.JSON(http.StatusOK, rd)

}

// Delete ALL users
func (base *Controller) DeleteUsers(c *gin.Context){
	err := admin.DeleteUsers()
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Unable to delete all users", err,nil)
		c.JSON(http.StatusInternalServerError,rd)
		return
	}
	rd := utility.BuildSuccessResponse(http.StatusOK, "success", gin.H{"message" : "All users deleted successfully"} )
	c.JSON(http.StatusOK, rd)

}



// Delete User
func (base *Controller) DeleteUser(c *gin.Context){
	username := c.Param("username") //string
	err := admin.DeleteUser(username)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "Unable to delete user", err,nil)
		c.JSON(http.StatusInternalServerError,rd)
		return
	}
	rd := utility.BuildSuccessResponse(http.StatusOK, "success", gin.H{"message" : "user deleted successfully"} )
	c.JSON(http.StatusOK, rd)
}

// this returns the mined images of all users
func (base *Controller) GetAllMinedImages(c *gin.Context) {

	secretKey := config.GetConfig().Server.Secret
	token := utility.ExtractToken(c)

	_, err := utility.GetKey("id", token, secretKey)

	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, "failed", "could not verify token", nil, gin.H{"error": err.Error()})
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

	minedImages, err := admin.GetMinedImages()
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", err.Error(), err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return

	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "success", gin.H{"data": minedImages})
	c.JSON(http.StatusOK, rd)

}

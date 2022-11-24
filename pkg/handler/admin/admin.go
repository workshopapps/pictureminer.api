package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/workshopapps/pictureminer.api/service/admin"
	"github.com/workshopapps/pictureminer.api/utility"
)

type Controller struct {
	Validate *validator.Validate
	Logger   *utility.Logger
}

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

// this returns the mined images of all users
func (base *Controller) GetAllMinedImages(c *gin.Context) {

	minedImages, err := admin.GetMinedImages()
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", err.Error(), err, nil)
		c.JSON(http.StatusInternalServerError, rd)
		return

	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "success", gin.H{"data": minedImages})
	c.JSON(http.StatusOK, rd)

}

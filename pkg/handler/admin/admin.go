package admin

import (
	"net/http"

	// "encoding/json"
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

	rd := utility.BuildSuccessResponse(http.StatusOK, "success", map[string]interface{}{"data": users})
	c.JSON(http.StatusOK, rd)

}

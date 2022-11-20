package health

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/service/ping"
	"github.com/workshopapps/pictureminer.api/utility"
)

type Controller struct {
	Validate *validator.Validate
	Logger   *utility.Logger
}

// Post             godoc
// @Summary      Checks the status of the server
// @Description  Send a dummy post request to test the status of the server
// @Tags         health
// @Produce      json
// @Param        ping  body      model.Ping  true  "Ping JSON"
// @Success      200   {object}  utility.Response
// @Router       /api/v1/health [post]
func (base *Controller) Post(c *gin.Context) {
	var (
		req = model.Ping{}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validate.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validate), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	if !ping.ReturnTrue() {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "ping failed", fmt.Errorf("ping failed"), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	base.Logger.Info("ping successfull")

	rd := utility.BuildSuccessResponse(http.StatusOK, "ping successfull", req.Message)
	c.JSON(http.StatusOK, rd)

}

// Get             godoc
// @Summary      Checks the status of the server
// @Description  Responds with the server status as JSON.
// @Tags         health
// @Produce      json
// @Success      200  {object}  utility.Response
// @Router       /api/v1/health [get]
func (base *Controller) Get(c *gin.Context) {
	if !ping.ReturnTrue() {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "ping failed", fmt.Errorf("ping failed"), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	base.Logger.Info("ping successfull")
	rd := utility.BuildSuccessResponse(http.StatusOK, "ping successfull", gin.H{"user": "user object"})
	c.JSON(http.StatusOK, rd)

}

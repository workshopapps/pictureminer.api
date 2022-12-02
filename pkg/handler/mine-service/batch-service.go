package mineservice

import (
  "fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/workshopapps/pictureminer.api/internal/config"
	mineservice "github.com/workshopapps/pictureminer.api/service/mine-service"
	"github.com/workshopapps/pictureminer.api/utility"
)



func (base *Controller) GetBatchResult(c *gin.Context) {

	secretKey := config.GetConfig().Server.Secret
	token := utility.ExtractToken(c)
	userId, err := utility.GetKey("id", token, secretKey)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, "failed", "could not verify token", nil, gin.H{"error": err.Error()})
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

  batchId := c.Param("batch_id")

  UserIdstr := fmt.Sprintf("%v", userId)

	batchImages, err := mineservice.GetbatchImages(UserIdstr,batchId)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "could get mined images", nil, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	c.JSON(http.StatusOK, batchImages)

}

package mineservice

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/workshopapps/pictureminer.api/internal/config"
	mineservice "github.com/workshopapps/pictureminer.api/service/mine-service"
	"github.com/workshopapps/pictureminer.api/utility"
)

func (base *Controller) ProcessBatch(c *gin.Context) {
	secretKey := config.GetConfig().Server.Secret
	token := utility.ExtractToken(c)
	userID, err := utility.GetKey("id", token, secretKey)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, "failed", "unable to verify token", gin.H{"error": err.Error()}, nil)
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

	if c.ContentType() != "multipart/form-data" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "invalid request", gin.H{"error": "file is not present"}, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	file, fileHeader, err := c.Request.FormFile("csv")
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "unable to process file", gin.H{"error": err.Error()}, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	defer file.Close()

	if strings.ToLower(filepath.Ext(fileHeader.Filename)) != ".csv" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "invalid file", gin.H{"error": "file must be a csv"}, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	// validate user ID
	id, ok := userID.(string)
	if !ok {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "invalid user id claim", gin.H{"error": "could not process user id"}, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	res, code, err := mineservice.ProcessBatchService(id, file)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "failed", "an error occurred", gin.H{"error": err.Error()}, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "process batch started", res)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) GetBatches(c *gin.Context) {

	secretKey := config.GetConfig().Server.Secret
	token := utility.ExtractToken(c)
	userID, err := utility.GetKey("id", token, secretKey)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, "failed", "could not verify token", nil, gin.H{"error": err.Error()})
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

	// validate user ID
	id, ok := userID.(string)
	if !ok {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "invalid user id claim", gin.H{"error": "could not process user id"}, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	batches, err := mineservice.GetBatchesService(id)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "could retrive batches", gin.H{"error": err.Error()}, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	c.JSON(http.StatusOK, batches)

}
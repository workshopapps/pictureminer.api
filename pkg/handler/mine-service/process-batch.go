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

	// retrieve and validate required create batch details
	batchName := c.PostForm("name")
	if batchName == "" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "invalid request", gin.H{"error": "name field missing"}, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	tags := c.PostForm("tags")
	if tags == "" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "invalid request", gin.H{"error": "tags field missing"}, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	// optional field
	desc := c.PostForm("description")

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

	res, code, err := mineservice.ProcessBatchService(id, batchName, desc, tags, file)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "failed", "an error occurred", gin.H{"error": err.Error()}, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "process batch started", res)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) ProcessBatchCSV(c *gin.Context) {
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

	res, code, err := mineservice.ProcessBatchCSVService(id, file)
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

func (base *Controller) GetBatchImages(c *gin.Context) {

	secretKey := config.GetConfig().Server.Secret
	token := utility.ExtractToken(c)
	_, err := utility.GetKey("id", token, secretKey)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, "failed", "could not verify token", nil, gin.H{"error": err.Error()})
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

	// get batch id
	batchID := c.Param("batch_id")
	if batchID == "" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "invalid request", gin.H{"error": "batch id field missing"}, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	batchImgs, err := mineservice.GetBatchImages(batchID)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "could not retrive batch images", gin.H{"error": err.Error()}, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	c.JSON(http.StatusOK, batchImgs)

}

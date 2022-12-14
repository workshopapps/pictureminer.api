package batch

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/workshopapps/pictureminer.api/internal/config"
	batchservice "github.com/workshopapps/pictureminer.api/service/batch-service"
	mineservice "github.com/workshopapps/pictureminer.api/service/mine-service"
	"github.com/workshopapps/pictureminer.api/utility"
)

type Controller struct {
	Validate *validator.Validate
	Logger   *utility.Logger
}

// ProcessBatchAPI godoc
// @Summary      Processes a batch of images
// @Description  Process a list of images as a batch
// @Tags         batch-api
// @Param       json formData file true "json"
// @Param       csv formData file false "csv"
// @Success      200   {object}  utility.Response
// @Router       /batch-service/process-batch-api [post]
// @Security BearerAuth
func (base *Controller) ProcessBatchAPI(c *gin.Context) {
	secretKey := config.GetConfig().Server.Secret
	token := utility.ExtractToken(c)
	userId, err := utility.GetKey("id", token, secretKey)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, "failed", "could not verify token", nil, gin.H{"error": err.Error()})
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

	UserIdstr := fmt.Sprintf("%v", userId)

	count, _ := mineservice.GetMonthlylimit(UserIdstr)
	if count != true {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Monthly limit exceeded", nil , gin.H{"error":" you have exceeded monthly limit of 10 Mine requests"})
		c.JSON(http.StatusBadRequest, rd)
		return
	}


	if c.ContentType() != "multipart/form-data" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "invalid request", nil, gin.H{"error": "file is not present"})
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	result, err := batchservice.ProcessBatchAPI(userId.(string), c.Request)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "undefined error", nil, err.Error())
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "process batch started", result)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) ProcessBatch(c *gin.Context) {
	secretKey := config.GetConfig().Server.Secret
	token := utility.ExtractToken(c)
	userID, err := utility.GetKey("id", token, secretKey)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, "failed", "unable to verify token", gin.H{"error": err.Error()}, nil)
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

	UserIdstr := fmt.Sprintf("%v", userID)

	count, _ := mineservice.GetMonthlylimit(UserIdstr)
	if count != true {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Monthly limit exceeded", nil , gin.H{"error":" you have exceeded monthly limit of 10 Mine requests"})
		c.JSON(http.StatusBadRequest, rd)
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

	res, code, err := batchservice.ProcessBatchService(id, batchName, desc, tags, file)
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

	UserIdstr := fmt.Sprintf("%v", userID)

	count, _ := mineservice.GetMonthlylimit(UserIdstr)
	if count != true {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Monthly limit exceeded", nil , gin.H{"error":" you have exceeded monthly limit of 10 Mine requests"})
		c.JSON(http.StatusBadRequest, rd)
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

	res, code, err := batchservice.ProcessBatchCSVService(id, file)
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

	batches, err := batchservice.GetBatchesService(id)
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

	resp, err := batchservice.GetBatchImages(batchID)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "could not retrive batch images", gin.H{"error": err.Error()}, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	c.JSON(http.StatusOK, resp)

}

func (base *Controller) DeleteBatch(c *gin.Context) {

	secretKey := config.GetConfig().Server.Secret
	token := utility.ExtractToken(c)
	_, err := utility.GetKey("id", token, secretKey)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, "failed", "could not verify token", nil, gin.H{"error": err.Error()})
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

	batchID := c.Param("id")

	err = batchservice.DeleteBatchService(batchID)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "could not delete batch", nil, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "delete batch success", gin.H{})
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) DownloadCsv(c *gin.Context) {

	// secretKey := config.GetConfig().Server.Secret
	// token := utility.ExtractToken(c)
	// userId, err := utility.GetKey("id", token, secretKey)
	// if err != nil {
	// 	rd := utility.BuildErrorResponse(http.StatusUnauthorized, "failed", "could not verify token", nil, gin.H{"error": err.Error()})
	// 	c.JSON(http.StatusUnauthorized, rd)
	// 	return
	// }
	batchId := c.Param("batchid")
	var dummySlice, err = batchservice.GetImagesInBatch(batchId)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "could not get images for this batch id", nil, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	errr := batchservice.ParseImageResponseForDownload(dummySlice)
	if errr != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "could not generate csv for download", nil, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	c.FileAttachment("filename.csv", batchId+".csv")
	defer os.Remove("filename.csv")

}

func (base *Controller) CountBatches(c *gin.Context) {
	secretKey := config.GetConfig().Server.Secret
	token := utility.ExtractToken(c)
	userId, err := utility.GetKey("id", token, secretKey)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, "failed", "could not verify token", gin.H{"error": err.Error()}, nil)
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

	userID, ok := userId.(string)
	if !ok {
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, "failed", "could not verify token", gin.H{"error": "invalid token type"}, nil)
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

	resp, code, err := batchservice.CountBatchesService(userID)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "failed", "could not retrieve batches counts", gin.H{"error": err.Error()}, nil)
		c.JSON(code, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "get batches count success", resp)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) CountProcess(c *gin.Context) {

	secretKey := config.GetConfig().Server.Secret
	token := utility.ExtractToken(c)
	userId, err := utility.GetKey("id", token, secretKey)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, "failed", "could not verify token", nil, gin.H{"error": err.Error()})
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

	UserIdstr := fmt.Sprintf("%v", userId)


	processCount, err := mineservice.ProcessCount(UserIdstr)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "could not get count", nil, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	c.JSON(http.StatusOK, processCount)
}

package mineservice

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/model"
	mineservice "github.com/workshopapps/pictureminer.api/service/mine-service"
	batchservice "github.com/workshopapps/pictureminer.api/service/batch-service"
	"github.com/workshopapps/pictureminer.api/utility"
)

type Controller struct {
	Validate *validator.Validate
	Logger   *utility.Logger
}

func (base *Controller) DemoMineImage(c *gin.Context) {
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

	minedImage, err := mineservice.DemoMineImage(image, imageHeader.Filename)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "undefined error", nil, err.Error())
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "mine image successful", minedImage)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) MineImageUpload(c *gin.Context) {

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

	minedImage, err := mineservice.MineServiceUpload(userId, image, imageHeader.Filename)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "undefined error", nil, err.Error())
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, "mine image successful", minedImage)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) MineImageUrl(c *gin.Context) {

	secretKey := config.GetConfig().Server.Secret
	token := utility.ExtractToken(c)
	userId, err := utility.GetKey("id", token, secretKey)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, "failed", "could not verify token", nil, gin.H{"error": err.Error()})
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

	var req model.MineImageUrlRequest

	err = c.Bind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "Unable to bind url parameter", nil, err.Error())
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validate.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "invalid url", nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	response, err := http.Get(req.Url)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "could not fetch image from url", nil, err.Error())
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	var image = response.Body
	defer image.Close()

	filename := getFileName(req.Url)

	if !utility.ValidImageFormat(filename) {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "invalid file", nil, gin.H{"error": "file is not an image"})
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	minedImage, err := mineservice.MineServiceUpload(userId, image, filename)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusServiceUnavailable, "failed", "could not save image", nil, err.Error())
		c.JSON(http.StatusServiceUnavailable, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, "image successfully mined", minedImage)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) GetMinedImages(c *gin.Context) {

	secretKey := config.GetConfig().Server.Secret
	token := utility.ExtractToken(c)
	userId, err := utility.GetKey("id", token, secretKey)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, "failed", "could not verify token", nil, gin.H{"error": err.Error()})
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

	minedImages, err := mineservice.GetMinedImages(userId)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "could get mined images", nil, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	c.JSON(http.StatusOK, minedImages)

}

func getFileName(url string) string {
	return url[strings.LastIndex(url, "/")+1:]
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
	errr := mineservice.ParseImageResponseForDownload(dummySlice)
	if errr != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "could not generate csv for download", nil, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	c.File("filename.csv")
	defer os.Remove("filename.csv")
	

}
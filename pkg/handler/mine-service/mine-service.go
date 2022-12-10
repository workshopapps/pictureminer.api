package mineservice

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/model"
	mineservice "github.com/workshopapps/pictureminer.api/service/mine-service"
	"github.com/workshopapps/pictureminer.api/service/user"
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

// Post             godoc
// @Summary     Mines an uploaded image
// @Description Send a post request containing a file an receives a response of its context content.
// @Tags        Mine-Service
// @Param       image formData file true "image"
// @Success     200  {object} utility.Response
// @Router      /mine-service/upload [post]
// @Security BearerAuth
func (base *Controller) MineImageUpload(c *gin.Context) {

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

	UserIdstr := fmt.Sprintf("%v", userId)

	count, _ := mineservice.GetMonthlylimit(UserIdstr)
	if count != true {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Monthly limit exceeded", nil , gin.H{"error":" you have exceeded monthly limit of 10 Mine requests"})
		c.JSON(http.StatusBadRequest, rd)
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

func (base *Controller) DeleteMinedImage(c *gin.Context) {

	secretKey := config.GetConfig().Server.Secret
	token := utility.ExtractToken(c)
	_, err := utility.GetKey("id", token, secretKey)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, "failed", "could not verify token", nil, gin.H{"error": err.Error()})
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

	// get image key
	imageKey := c.Param("key")
	if imageKey == "" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "invalid request", gin.H{"error": "key field missing"}, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = mineservice.DeleteMinedImageService(imageKey)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "could not delete mined image", nil, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "delete mined image success", gin.H{})
	c.JSON(http.StatusOK, rd)

}

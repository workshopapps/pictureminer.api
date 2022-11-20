package mineservice

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/workshopapps/pictureminer.api/internal/model"
	mineservice "github.com/workshopapps/pictureminer.api/service/mine-service"
	"github.com/workshopapps/pictureminer.api/utility"
)

type Controller struct {
	Validate *validator.Validate
	Logger   *utility.Logger
}

func (base *Controller) Post(c *gin.Context) {

	// TODO:Authorize request via JWT token

	if c.ContentType() != "multipart/form-data" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "invalid request", nil, gin.H{"error": "file is not present"})
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	image, fh, err := c.Request.FormFile("image")
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "couldn't parse file", nil, gin.H{"error": "image size is too large, must be less than 1MB"})
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	filename := fh.Filename
	if !strings.HasSuffix(filename, ".png") && !strings.HasSuffix(filename, ".jpg") && !strings.HasSuffix(filename, ".jpeg") {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "invalid file", nil, gin.H{"error": "file is not an image"})
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	minedImage, err := mineservice.MineServiceUpload(image, filename)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "server error", nil, err.Error())
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, "image successfully mined", minedImage)
	c.JSON(http.StatusOK, rd)
}

func (base *Controller) MineImageUrl(c *gin.Context) {

	req := model.MineImageUrlRequest{}

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "internal server error", nil, err.Error())
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

	defer response.Body.Close()

	urlSplit := strings.Split(req.Url, "/")[1:]
	urlSlice := urlSplit[len(urlSplit)-1:]
	var filename string = urlSlice[0]

	if !strings.HasSuffix(filename, ".png") && !strings.HasSuffix(filename, ".jpg") && !strings.HasSuffix(filename, ".jpeg") {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "invalid file", nil, gin.H{"error": "file is not an image"})
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	minedImage, err := mineservice.MineServiceUpload(image, filename)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "could not save image", nil, err.Error())
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, "image successfully mined", minedImage)
	c.JSON(http.StatusOK, rd)
}

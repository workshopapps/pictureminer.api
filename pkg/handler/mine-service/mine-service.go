package mineservice

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

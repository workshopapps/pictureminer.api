package mineservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/model"
	mineservice "github.com/workshopapps/pictureminer.api/service/mine-service"
	"github.com/workshopapps/pictureminer.api/utility"
)

func (base *Controller) PromptMineImageUpload(c *gin.Context) {

	secretKey := config.GetConfig().Server.Secret
	token := utility.ExtractToken(c)
	_, err := utility.GetKey("id", token, secretKey)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, "failed", "could not verify token", nil, gin.H{"error": err.Error()})
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

	var promptReq model.MineImagePromptRequest
	if err := c.BindQuery(&promptReq); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "invalid request", nil, gin.H{"error": "prompt not specified"})
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if err := base.Validate.Struct(&promptReq); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "invalid request", nil, err.Error())
		c.JSON(http.StatusBadRequest, rd)
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

	if !validImageFormat(imageHeader.Filename) {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "invalid file", nil, gin.H{"error": "file is not an image"})
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	minedImage, err := mineservice.PromptMineServiceUpload(image, imageHeader.Filename, promptReq.Prompt)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "undefined error", nil, err.Error())
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, "mine image successful", minedImage)
	c.JSON(http.StatusOK, rd)

}

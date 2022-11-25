package mineservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/model"
	mineservice "github.com/workshopapps/pictureminer.api/service/mine-service"
	"github.com/workshopapps/pictureminer.api/utility"
)

const FailedStatus = `failed`
const (
	VerifiedTokenMessage             = `Could not verify token`
	InvalidRequestMessage            = `Invalid request`
	InvalidFileMessage               = `Invalid file`
	InvalidUrlMessage                = `Invalid url`
	SuccessfulMiningMessage          = `Mine image successful`
	UndefinedErrorMessage            = `Undefined error`
	UnableToFetchImageFromUrlMessage = "Could not fetch image from url"
	UnableToParseFileMessage         = "Could not parse file"
	UnableToBindUrlMessage           = "Unable to bind url parameter"
)

const (
	PromptNotSpecifiedError = "Prompt Param not specified"
	FileIsNotAnImageError   = "File is not an image"
	FileIsNotPresentError   = "File is not present"
)

func (base *Controller) PromptMineImageUpload(c *gin.Context) {
	err := isTokenVerified(c)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, FailedStatus, VerifiedTokenMessage, nil, makeErrorMap(err.Error()))
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

	prompt := c.Query("prompt")
	if prompt == "" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, FailedStatus, InvalidRequestMessage, nil, makeErrorMap(PromptNotSpecifiedError))
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if c.ContentType() != "multipart/form-data" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, FailedStatus, InvalidRequestMessage, nil, makeErrorMap(FileIsNotPresentError))
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	image, imageHeader, err := c.Request.FormFile("image")
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, FailedStatus, UnableToParseFileMessage, nil, makeErrorMap(err.Error()))
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	defer image.Close()

	if !validImageFormat(imageHeader.Filename) {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, FailedStatus, InvalidFileMessage, nil, makeErrorMap(FileIsNotAnImageError))
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	minedImage, err := mineservice.PromptMineServiceUpload(image, imageHeader.Filename, prompt)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, FailedStatus, UndefinedErrorMessage, nil, err.Error())
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, SuccessfulMiningMessage, minedImage)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) PromptMineImageUrl(c *gin.Context) {

	err := isTokenVerified(c)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, FailedStatus, VerifiedTokenMessage, nil, makeErrorMap(err.Error()))
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

	var req model.MineImageUrlRequest

	err = c.Bind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, FailedStatus, UnableToBindUrlMessage, nil, err.Error())
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validate.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, FailedStatus, InvalidUrlMessage, nil, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	prompt := c.Query("prompt")
	if prompt == "" {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, FailedStatus, InvalidRequestMessage, nil, makeErrorMap(PromptNotSpecifiedError))
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	response, err := http.Get(req.Url)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, FailedStatus, UnableToFetchImageFromUrlMessage, nil, err.Error())
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	var image = response.Body
	defer image.Close()

	filename := getFileName(req.Url)

	if !validImageFormat(filename) {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, FailedStatus, InvalidFileMessage, nil, makeErrorMap(FileIsNotAnImageError))
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	minedImage, err := mineservice.PromptMineServiceUpload(image, filename, prompt)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, FailedStatus, UndefinedErrorMessage, nil, err.Error())
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusCreated, SuccessfulMiningMessage, minedImage)
	c.JSON(http.StatusOK, rd)

}

func isTokenVerified(c *gin.Context) (err error) {
	secretKey := config.GetConfig().Server.Secret
	token := utility.ExtractToken(c)
	_, err = utility.GetKey("id", token, secretKey)
	return err
}

func makeErrorMap(errorMesage string) gin.H {
	errorMap := gin.H{"error": errorMesage}
	return errorMap
}

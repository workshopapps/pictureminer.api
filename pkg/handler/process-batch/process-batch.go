package batch

import (
	"encoding/json"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/model"
	mineservice "github.com/workshopapps/pictureminer.api/service/mine-service"
	"github.com/workshopapps/pictureminer.api/utility"
)

type Controller struct {
	Validate *validator.Validate
	Logger   *utility.Logger
}

func (base *Controller) ProcessBatch(c *gin.Context) {
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

	jsonReq, csvFile, err := processJSONAndCSV(c.Request)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "undefined error", nil, err.Error())
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	result, err := mineservice.ProcessBatch(userId.(string), jsonReq, csvFile)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "undefined error", nil, err.Error())
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "mine image successful", result)
	c.JSON(http.StatusOK, rd)
}

// processJSONFile processes the JSON file to get the json data and determines if the JSON file
// contains either image URLs or a CSV file is passed in with the request
func processJSONAndCSV(r *http.Request) (*model.ProcessBatchRequest, multipart.File, error) {
	var jsonReq model.ProcessBatchRequest

	file, _, err := r.FormFile("json")
	if err != nil {
		return nil, nil, err
	}

	if err = json.NewDecoder(file).Decode(&jsonReq); err != nil {
		return nil, nil, err
	}

	if len(jsonReq.Images) != 0 {
		return &jsonReq, nil, nil
	}

	csvFile, _, err := r.FormFile("csv")
	if err != nil {
		return nil, nil, err
	}

	return &jsonReq, csvFile, nil
}

package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	mineservice "github.com/workshopapps/pictureminer.api/pkg/handler/mine-service"
	"github.com/workshopapps/pictureminer.api/utility"
)

func MineService(r *gin.Engine, validate *validator.Validate, ApiVersion string, logger *utility.Logger) *gin.Engine {
	mineservice := mineservice.Controller{Validate: validate, Logger: logger}

	authUrl := r.Group(fmt.Sprintf("/api/%v", ApiVersion))
	{
		authUrl.POST("/mine-service/upload", mineservice.MineImageUpload)
		authUrl.POST("/mine-service/url", mineservice.MineImageUrl)
		authUrl.GET("/mine-service/get-all", mineservice.GetMinedImages)
		authUrl.GET("/batch-service/get-all/:batch_id", mineservice.GetBatchResult)
		authUrl.POST("/mine-service/demo", mineservice.DemoMineImage)
	}
	return r
}

func DownLoadCsv(r *gin.Engine, validate *validator.Validate, ApiVersion string, logger *utility.Logger) *gin.Engine {
	mineservice := mineservice.Controller{Validate: validate, Logger: logger}
	authUrl := r.Group(fmt.Sprintf("/api/%v", ApiVersion))
	{
		authUrl.GET("/csv-service/download/:batchid", mineservice.DownloadCsv)
	}

	return r
}

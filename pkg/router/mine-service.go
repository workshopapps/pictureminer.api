package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	mineservice "github.com/workshopapps/pictureminer.api/pkg/handler/mine-service"
	"github.com/workshopapps/pictureminer.api/utility"
)

func MineServiceUpload(r *gin.Engine, validate *validator.Validate, ApiVersion string, logger *utility.Logger) *gin.Engine {
	mineservice := mineservice.Controller{Validate: validate, Logger: logger}

	authUrl := r.Group(fmt.Sprintf("/api/%v", ApiVersion))
	{
		authUrl.POST("/mine-service", mineservice.Post)
		authUrl.POST("/mine-service/url", mineservice.MineImageUrl)
	}
	return r
}

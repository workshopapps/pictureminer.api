package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/workshopapps/pictureminer.api/pkg/handler/user"
	"github.com/workshopapps/pictureminer.api/utility"
)

func Signup(r *gin.Engine, validate *validator.Validate, ApiVersion string, logger *utility.Logger) *gin.Engine {
	//add
	User := user.Controller{Validate: validate, Logger: logger}

	//authUrl := r.Group(fmt.Sprintf("/api/%v", ApiVersion))
	//{
	//	authUrl.POST("/health", health.Post)
	//	authUrl.GET("/health", health.Get)
	//}
	r.POST("/"+"signup", User.Signup)

	return r
}

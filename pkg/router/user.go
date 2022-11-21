package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/workshopapps/pictureminer.api/pkg/handler/user"
	"github.com/workshopapps/pictureminer.api/utility"
)

func Auth(r *gin.Engine, validate *validator.Validate, ApiVersion string, logger *utility.Logger) *gin.Engine {

	auth := user.Controller{Validate: validate, Logger: logger}

	authUrl := r.Group(fmt.Sprintf("/api/%v", ApiVersion))
	{
		authUrl.POST("/create_user", auth.CreateUser)
		authUrl.POST("/login", auth.Login)
		authUrl.POST("/signup", auth.Signup)
	}
	return r
}

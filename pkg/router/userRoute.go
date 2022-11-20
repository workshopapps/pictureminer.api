package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/workshopapps/pictureminer.api/pkg/handler/user"
	"github.com/workshopapps/pictureminer.api/utility"
)

func Signup(r *gin.Engine, validate *validator.Validate, ApiVersion string, logger *utility.Logger) *gin.Engine {
	//add
	user := user.Controller{Validate: validate, Logger: logger}

	r.POST(fmt.Sprintf("/api/%v", ApiVersion) + "/signup", user.Signup)

	return r
}

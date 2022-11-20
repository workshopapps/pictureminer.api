package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/workshopapps/pictureminer.api/pkg/handler/admin"
  // "github.com/workshopapps/pictureminer.api/pkg/handler/user"
	"github.com/workshopapps/pictureminer.api/utility"
)

func Admin(r *gin.Engine, validate *validator.Validate, ApiVersion string, logger *utility.Logger) *gin.Engine {
	admin := admin.Controller{Validate: validate, Logger: logger}

	adminUrl := r.Group(fmt.Sprintf("/api/%v", ApiVersion))
	{
		adminUrl.GET("/admin/users", admin.GetUsers)
	}
	return r
}

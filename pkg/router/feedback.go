package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/workshopapps/pictureminer.api/pkg/handler/feedback"
	"github.com/workshopapps/pictureminer.api/utility"
)

func Feedback(r *gin.Engine, validate *validator.Validate, ApiVersion string, logger *utility.Logger) *gin.Engine {

	feedback:= feedback.Controller{Validate: validate, Logger: logger}
	feedbackUrl := r.Group(fmt.Sprintf("/api/%v", ApiVersion))
	{
		feedbackUrl.POST("/feedback", feedback.AcceptFeedback)
		feedbackUrl.GET("/feedback/all", feedback.GetAllFeedback)
	}
	return r
}
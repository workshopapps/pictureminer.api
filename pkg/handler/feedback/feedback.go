package feedback

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/service/feedback"
	"github.com/workshopapps/pictureminer.api/utility"
)

type Controller struct {
	Validate *validator.Validate
	Logger   *utility.Logger
}

func (base *Controller) AcceptFeedback(c *gin.Context){
	var Feedback model.Feedback

	err:= c.Bind(&Feedback)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Unable to bind user signup details", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	feedback, err := feedback.CollectFeedbackFromUser(Feedback)
	if err != nil {
		rd := utility.BuildErrorResponse(400, "error", "could not parse feedback", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	c.JSON(http.StatusCreated, feedback)
}

func (base *Controller) GetAllFeedback(c *gin.Context){
	reviews, err := feedback.GetAllFeedback()

	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "could not get all the reviews", nil, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	c.JSON(http.StatusOK, reviews)
}

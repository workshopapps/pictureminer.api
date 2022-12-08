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

// Post             godoc
// @Summary    Sends feedback to discripto
// @Description Send feedback to discripto
// @Tags        Feedback
// @Param       Feedback body model.FeedbackRequest true "Create feedback" model.Feedback
// @Success     200  {object} model.FeedbackCreatedResponse
// @Failure  	400 {object} utility.Response
// @Router      /feedback [post]
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

// Get             godoc
// @Summary     Gets all feedback sent discripto
// @Description Gets all feedback sent to discripto
// @Tags        Feedback
// @Success     200  {object} []model.Feedback
// @Failure  	400 {object} utility.Response
// @Router      /feedback/all [get]
func (base *Controller) GetAllFeedback(c *gin.Context){
	reviews, err := feedback.GetAllFeedback()

	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "could not get all the reviews", nil, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	c.JSON(http.StatusOK, reviews)
}

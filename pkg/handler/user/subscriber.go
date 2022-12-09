package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/service/user"
	"github.com/workshopapps/pictureminer.api/utility"
)


func (base *Controller) SubscriberEmail(c *gin.Context) {

	// bind emails to SubscriberEmail struct
	var subscriberEmail model.SubscriberEmail
	err := c.Bind(&subscriberEmail)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Unable to bind subscriber email details", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validate.Struct(&subscriberEmail)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validate), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	SubscriberEmailResponse, msg, code, err := user.SubscriberEmailResponse(subscriberEmail)
	if err != nil {
		rd := utility.BuildErrorResponse(code, "error", msg, err, nil)
		c.JSON(code, rd)
		return
	}

	object := utility.BuildSuccessResponse(200, "Email Submission successful", SubscriberEmailResponse)
	c.JSON(200, object)
}

func (base *Controller) GetSubscribtion(c *gin.Context) {

	// bind emails to SubscriberEmail struct
	email, ok:= c.GetQuery("user")
	if !ok {
		rd := utility.BuildErrorResponse(400, "error", "please supply the user email", "", nil)
		c.JSON(400, rd)
		return
	}

	sub, err:= user.GetUserSubscribtion(email)
	if err != nil {
		rd := utility.BuildErrorResponse(400, "error", "USer may not have subscribtion", err, nil)
		c.JSON(400, rd)
		return
	}

	object := utility.BuildSuccessResponse(200, "Email Submission successful", sub)
	c.JSON(200, object)
}

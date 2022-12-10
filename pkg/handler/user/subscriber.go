package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/service/user"
	"github.com/workshopapps/pictureminer.api/utility"
)

// Create Subscription godoc
// @Summary      Create Subscription
// @Description  create a user subscription
// @Tags         users
// @Produce      json
// @Param User body model.SubscriptionRequest true "User Subscription" model.SubscriberEmail
// @Success      200  {object}   model.SubscriberEmail
// @Failure      400  {object}  utility.Response
// @Router       /subscription [post]
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

	object := utility.BuildSuccessResponse(200, "Email Subcription successful", SubscriberEmailResponse)
	c.JSON(200, object)
}

// Fetch Subscription godoc
// @Summary      Fetch Subscription
// @Description  get subscription status information
// @Tags         users
// @Produce      json
// @Param        user    query     string  false  "email of the user"  Format(email)
// @Success      200  {object}   model.SubscriberEmail
// @Failure      400  {object}  utility.Response
// @Router       /subscription [get]
func (base *Controller) GetSubscription(c *gin.Context) {

	// bind emails to SubscriberEmail struct
	email, ok := c.GetQuery("user")
	if !ok {
		rd := utility.BuildErrorResponse(400, "error", "please supply the user email", "", nil)
		c.JSON(400, rd)
		return
	}

	sub, err := user.GetUserSubscription(email)
	if err != nil {
		rd := utility.BuildErrorResponse(400, "error", "USer may not have subscription", err, nil)
		c.JSON(400, rd)
		return
	}

	object := utility.BuildSuccessResponse(200, "User subscription retrieved sucessfully", sub)
	c.JSON(200, object)
}

func (base *Controller) VerifyEmail(c *gin.Context) {
	secretKey := config.GetConfig().Server.Secret
	token := utility.ExtractToken(c)
	userId, err := utility.GetKey("id", token, secretKey)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnauthorized, "failed", "could not verify token", nil, gin.H{"error": err.Error()})
		c.JSON(http.StatusUnauthorized, rd)
		return
	}

	if err := user.VerifyEmail(userId); err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "failed", "could not verify token", nil, gin.H{"error": err.Error()})
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	rd := utility.BuildSuccessResponse(http.StatusAccepted, "email successfully verified", nil)
	c.JSON(http.StatusOK, rd)
}

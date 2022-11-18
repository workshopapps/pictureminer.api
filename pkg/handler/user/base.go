package user

import (
	"github.com/go-playground/validator/v10"
	"github.com/workshopapps/pictureminer.api/utility"
)

type Controller struct {
	Validate *validator.Validate
	Logger   *utility.Logger
}

type UserLoginField struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

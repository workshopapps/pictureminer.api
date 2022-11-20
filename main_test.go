package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	us "github.com/workshopapps/pictureminer.api/pkg/handler/user"
	os "github.com/workshopapps/pictureminer.api/pkg/handler/health"
	"github.com/workshopapps/pictureminer.api/utility"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

type Controller struct {
	Validate *validator.Validate
	Logger   *utility.Logger
}

type User struct {
	Username  string `json:"username" validate:"required,min=2,max=100"`
	FirstName    string             `bson:"first_name" json:"first_name"`
	LastName     string             `bson:"last_name" json:"last_name"`
	Email     string `json:"email" validate:"email,required"`
	Password  string `json:"Password" validate:"required,min=6"`
}

var RANDOM = xid.New().String()

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}


func TestCreateUserHandler(t *testing.T) {
	var validate *validator.Validate
	var logger *utility.Logger
	Auth := us.Controller{Validate: validate, Logger: logger}

	r := SetUpRouter()
	r.POST("api/v1/create_user", Auth.CreateUser)
	var user User

	// random := xid.New().String()

	user.Username = "workshopapps" + RANDOM
	user.FirstName = "Christopher"
	user.LastName = "Nwokoye"
	user.Email = "workshopapps" + RANDOM + "@gmail.com"
	user.Password = "blockchain"

	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "api/v1/create_user", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	// fmt.Println(bytes.NewBuffer(jsonValue))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}



func TestLoginUserHandler(t *testing.T) {

	var validate *validator.Validate
	var logger *utility.Logger
	r := SetUpRouter()
	Auth := us.Controller{Validate: validate, Logger: logger}
	r.POST("api/v1/login", Auth.Login)

	var user User
	user.Email = "workshopapps" + RANDOM + "@gmail.com"
	user.Password = "blockchain"

	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "api/v1/login", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

package admin

import (
	"fmt"
	"net/http"
"context"
"time"
// "encoding/json"
"go.mongodb.org/mongo-driver/bson"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/go-playground/validator/v10"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"github.com/workshopapps/pictureminer.api/service/ping"
	"github.com/workshopapps/pictureminer.api/utility"
	"github.com/workshopapps/pictureminer.api/internal/constants"
)

type Controller struct {
	Validate *validator.Validate
	Logger   *utility.Logger
}

// mongoClient := mongodb.Connection()
var userCollection *mongo.Collection = mongodb.GetCollection(mongodb.Connection(), constants.UserDatabase, constants.UserCollection)
var validate = validator.New()
// var userCollection *mongo.Collection = database.GetCollection(database.Client,constants.UserDatabase , constants.UserCollection)
// var validate = validator.New()
//Database connection
// mongoClient := mongodb.Connection()
// userCollection := mongodb.GetCollection(mongoClient, constants.UserDatabase, constants.UserCollection)

func (base *Controller) GetUsers(c *gin.Context) {
var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)


	if !ping.ReturnTrue() {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "ping failed", fmt.Errorf("ping failed"), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	base.Logger.Info("ping successfull")
	rd := utility.BuildSuccessResponse(http.StatusOK, "ping successfull", gin.H{"user": "user object"})
	c.JSON(http.StatusOK, rd)


var users []model.User
	      defer cancel()

filter :=  bson.M{}

cursor, err := userCollection.Find(ctx, filter)
if err != nil {
	panic(err)
}
// end find

//reading from the db in an optimal way
			 defer cursor.Close(ctx)
			 for cursor.Next(ctx) {
					 var singleUser model.User
					 if err = cursor.Decode(&singleUser); err != nil {
						 rd = utility.BuildErrorResponse(http.StatusInternalServerError, "error", "ping failed" , fmt.Errorf("display failed"), nil)
						 // rd := utility.BuildErrorResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}

							 c.JSON(http.StatusInternalServerError,rd )
					 }

					 users = append(users, singleUser)
			 }
rd = utility.BuildSuccessResponse(http.StatusOK, "success", map[string]interface{}{"data": users})
			 c.JSON(http.StatusOK,rd,)

}

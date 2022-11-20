package admin

import (
	"fmt"
	"context"
	"time"
	"log"
	// "net/http"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/workshopapps/pictureminer.api/internal/model"
	// "github.com/workshopapps/pictureminer.api/service/ping"
	"github.com/workshopapps/pictureminer.api/utility"
	"github.com/workshopapps/pictureminer.api/internal/constants"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
)

type Controller struct {
	Validate *validator.Validate
	Logger   *utility.Logger
}


var userCollection *mongo.Collection = mongodb.GetCollection(mongodb.Client,constants.UserDatabase , constants.UserCollection)
var validate = validator.New()

func (base *Controller) GetUsers(c *gin.Context) {
	// if !ping.ReturnTrue() {
	// 	rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "ping failed", fmt.Errorf("ping failed"), nil)
	// 	c.JSON(http.StatusBadRequest, rd)
	// 	return
	// }
	// base.Logger.Info("ping successfull")
	// rd := utility.BuildSuccessResponse(http.StatusOK, "ping successfull", gin.H{"user": "user object"})
	// c.JSON(http.StatusOK, rd)

	// var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)


filter := bson.D{{"username", ""}}

cursor, err := userCollection.Find(ctx, filter)
if err != nil {
	// panic(err)
	log.Println(err)
}
// end find

var users []model.User
if err = cursor.All(ctx, &users); err != nil {
	// panic(err)
	log.Println(err)
}

for _, user := range users {
	cursor.Decode(&user)
	output, err := json.MarshalIndent(user, "", "    ")
	if err != nil {
		// panic(err)
		log.Println(err)
	}
	defer cancel()
	fmt.Printf("%s\n", output)
}

}

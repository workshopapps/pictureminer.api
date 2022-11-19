package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	// "github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/service/ping"
	"github.com/workshopapps/pictureminer.api/utility"
)

type Controller struct {
	Validate *validator.Validate
	Logger   *utility.Logger
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
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

// begin find
coll := client.Database("sample_restaurants").Collection("restaurants")
filter := bson.D{{"cuisine", "Italian"}}

cursor, err := coll.Find(context.TODO(), filter)
if err != nil {
	panic(err)
}
// end find

var results []Restaurant
if err = cursor.All(context.TODO(), &results); err != nil {
	panic(err)
}

for _, result := range results {
	cursor.Decode(&result)
	output, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", output)
}

}

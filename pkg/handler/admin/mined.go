package admin

import (
	"fmt"
	"net/http"
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"github.com/workshopapps/pictureminer.api/service/ping"
	"github.com/workshopapps/pictureminer.api/utility"
	"github.com/workshopapps/pictureminer.api/internal/constants"
)


var imageCollection *mongo.Collection = mongodb.GetCollection(mongodb.Connection(), constants.UserDatabase, constants.ImageCollection)

func (base *Controller) GetImages(c *gin.Context) {
var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

if !ping.ReturnTrue() {
		rd := utility.BuildErrorResponse(http.StatusInternalServerError, "error", "ping failed", fmt.Errorf("ping failed"), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}
	base.Logger.Info("ping successfull")
	rd := utility.BuildSuccessResponse(http.StatusOK, "ping successfull", gin.H{"user": "user object"})
	c.JSON(http.StatusOK, rd)


var images []model.MinedImage
	      defer cancel()

filter :=  bson.M{}

cursor, err := imageCollection.Find(ctx, filter)
if err != nil {
	panic(err)
}
// end find

//reading from the db in an optimal way
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
	 var singleImage model.MinedImage
	 if err = cursor.Decode(&singleImage); err != nil {
		 rd = utility.BuildErrorResponse(http.StatusInternalServerError, "error", "ping failed" , fmt.Errorf("display failed"), nil)

			 c.JSON(http.StatusInternalServerError,rd )
	 }

	 images = append(images, singleImage)
	 }
rd = utility.BuildSuccessResponse(http.StatusOK, "success", map[string]interface{}{"data": images})
			 c.JSON(http.StatusOK,rd,)

}

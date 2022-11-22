package admin

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/constants"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUsers() ([]model.User, error) {

	ctx := context.TODO()
	cursor, err := mongodb.SelectFromCollection(ctx, config.GetConfig().Mongodb.Database, constants.UserCollection, bson.M{})
	if err != nil {
		return []model.User{}, err
	}

	var users []model.User
	cursor.All(ctx, &users)

	return users, nil
}

// DeleteUser Function
func DeleteUser() gin.HandlerFunc{
	return func(c *gin.Context) {
		userId := c.Param("_id")

		var userCollection *mongo.Collection = mongodb.OpenCollection(mongodb.Client, "user")
		var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		var user *model.User

		res, err := userCollection.DeleteOne(ctx, bson.D{{userId, &user.User_Id}})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fmt.Printf("deleted %v documents\n", res.DeletedCount)
		c.JSON(http.StatusOK, &user)
	}
}

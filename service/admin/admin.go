package admin

import (
	"context"
	"time"

	"github.com/workshopapps/pictureminer.api/internal/constants"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUsers() ([]model.User, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var users []model.User
	filter := bson.M{}
	cursor, err := mongodb.SelectFromCollection(ctx, constants.UserDatabase, constants.UserCollection, filter)
	if err != nil {
		return []model.User{}, err
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var singleUser model.User
		if err = cursor.Decode(&singleUser); err != nil {
			continue
		} else {
			users = append(users, singleUser)
		}

	}
	return users, nil
}

package admin

import (
	"context"

	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/constants"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson"
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

func GetMinedImages() ([]model.MinedImage, error){

	ctx := context.TODO()
	cursor, err := mongodb.SelectFromCollection(ctx, config.GetConfig().Mongodb.Database, constants.ImageCollection, bson.M{})

	if err != nil {
		return []model.MinedImage{}, err
	}

	var minedImages []model.MinedImage
	cursor.All(ctx, &minedImages)

	return minedImages, nil
}

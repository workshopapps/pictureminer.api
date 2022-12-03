package mineservice

import (
	"context"

	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/constants"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

func GetBatchesService(userID string) ([]model.Batch, error) {
	db := config.GetConfig().Mongodb.Database
	ctx := context.Background()
	filter := bson.M{"user_id": userID}
	cursor, err := mongodb.SelectFromCollection(ctx, db, constants.BatchCollection, filter)
	if err != nil {
		return nil, err
	}

	batches := []model.Batch{}
	cursor.All(ctx, &batches)

	return batches, nil
}

package batchservice

import (
	"context"

	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/constants"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func GetBatchImages(batchID string) (interface{}, error) {
	db := config.GetConfig().Mongodb.Database
	ctx := context.Background()

	// check if id is valid
	validBatchID, err := primitive.ObjectIDFromHex(batchID)
	if err != nil {
		return []model.BatchImage{}, err
	}

	batch := model.Batch{}
	batchCollection := mongodb.GetCollection(mongodb.Connection(), db, constants.BatchCollection)
	batchCollection.FindOne(ctx, bson.M{"_id": validBatchID}).Decode(&batch)

	filter := bson.M{"batch_id": batchID}
	cursor, err := mongodb.SelectFromCollection(ctx, db, constants.BatchImageCollection, filter)
	if err != nil {
		return nil, err
	}

	imgs := []model.BatchImage{}
	err = cursor.All(ctx, &imgs)
	if err != nil {
		return []model.BatchImage{}, err
	}

	mapper := map[string][]string{}
	mapper[Untagged] = []string{}
	for _, tag := range batch.Tags {
		mapper[tag] = nil
	}

	for _, img := range imgs {
		mapper[img.Tag] = append(mapper[img.Tag], img.URL)
	}

	resp := map[string][]string{}
	for key, val := range mapper {
		if val != nil {
			resp[key] = val
		}
	}

	return resp, nil
}

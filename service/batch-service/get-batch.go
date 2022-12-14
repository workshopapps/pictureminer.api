package batchservice

import (
	"context"
	"net/http"

	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/constants"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetBatchesService(userID string) (interface{}, error) {
	db := config.GetConfig().Mongodb.Database
	ctx := context.Background()
	filter := bson.M{"user_id": userID}
	cursor, err := mongodb.SelectFromCollection(ctx, db, constants.BatchCollection, filter)
	if err != nil {
		return nil, err
	}

	batches := []model.BatchResponse{}
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
		return []model.BatchImage{}, err
	}

	imgs := []model.BatchImage{}
	err = cursor.All(ctx, &imgs)
	if err != nil {
		return []model.BatchImage{}, err
	}

	mapper := map[string][]string{}
	mapper[Untagged] = nil
	for _, tag := range batch.Tags {
		mapper[tag] = nil
	}

	for _, img := range imgs {
		mapper[img.Tag] = append(mapper[img.Tag], img.URL)
	}

	mapped := map[string][]string{}
	for key, val := range mapper {
		if val != nil {
			mapped[key] = val
		}
	}

	resp := []map[string][]string{}
	for key, val := range mapped {
		resp = append(resp, map[string][]string{key: val})
	}

	return resp, nil
}

func GetImagesInBatch(batchId string) ([]model.BatchImage, error) {
	db := config.GetConfig().Mongodb.Database
	ctx := context.Background()
	filter := bson.M{"batch_id": batchId}

	cursor, err := mongodb.SelectFromCollection(ctx, db, constants.BatchImageCollection, filter)
	if err != nil {
		return nil, err
	}

	batchImages := make([]model.BatchImage, 0)
	cursor.All(ctx, &batchImages)

	return batchImages, nil
}

func CountBatchesService(userID string) (interface{}, int, error) {
	db := config.GetConfig().Mongodb.Database
	ctx := context.Background()

	// get all user's  batches
	filter := bson.M{"user_id": userID}
	bcursor, err := mongodb.SelectFromCollection(ctx, db, constants.BatchCollection, filter)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer bcursor.Close(ctx)

	// for O(1) lookups
	batchesMap := map[string]model.Batch{}
	for bcursor.Next(ctx) {
		var b model.Batch
		bcursor.Decode(&b)
		batchesMap[b.ID.Hex()] = b
	}

	// get cursor for all images
	icursor, err := mongodb.SelectFromCollection(ctx, db, constants.BatchImageCollection, bson.M{})
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	respMap := map[string]model.BatchesCountResponse{}
	// decode and process N images at a time to prevent high memory usage
	N := 1_000_000
	for icursor.Next(ctx) {
		imgs := decodeNImages(ctx, icursor, N)
		for _, img := range imgs {
			if _, ok := batchesMap[img.BatchID]; !ok {
				continue
			}
			if _, ok := respMap[img.BatchID]; !ok {
				tagged, untagged := getDefaultTagsCount(img.Tag)
				cbr := model.BatchesCountResponse{
					ID:          img.BatchID,
					Untagged:    untagged,
					Tagged:      tagged,
					DateCreated: batchesMap[img.BatchID].DateCreated,
				}
				respMap[img.BatchID] = cbr
			} else {
				tagged, untagged := getDefaultTagsCount(img.Tag)
				cbr := respMap[img.BatchID]
				cbr.Tagged += tagged
				cbr.Untagged += untagged
				respMap[img.BatchID] = cbr
			}
		}
	}

	resp := []model.BatchesCountResponse{}
	for _, cbr := range respMap {
		resp = append(resp, cbr)
	}

	return resp, http.StatusOK, nil
}

func getDefaultTagsCount(tag string) (int, int) {
	untagged, tagged := 0, 0
	if tag == Untagged {
		untagged = 1
	} else {
		tagged = 1
	}
	return tagged, untagged
}

func decodeNImages(ctx context.Context, cursor *mongo.Cursor, N int) []model.BatchImage {
	imgs := []model.BatchImage{}
	for i := 0; i < N && cursor.Next(ctx); i++ {
		var img model.BatchImage
		cursor.Decode(&img)
		imgs = append(imgs, img)
	}
	return imgs
}

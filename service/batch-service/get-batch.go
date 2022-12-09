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
	filter := bson.M{"user_id": userID}

	// get cursor for all batches
	bcursor, err := mongodb.SelectFromCollection(ctx, db, constants.BatchCollection, filter)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer bcursor.Close(ctx)

	bImgCol := mongodb.GetCollection(mongodb.Connection(), db, constants.BatchImageCollection)
	taggedTotal, untaggedTotal := 0, 0

	for bcursor.Next(ctx) {
		var b model.Batch
		bcursor.Decode(&b)
		tagged, untagged := countImageTags(ctx, b, bImgCol)
		taggedTotal += tagged
		untaggedTotal += untagged
	}

	resp := model.BatchesCountResponse{
		Total:    taggedTotal + untaggedTotal,
		Tagged:   taggedTotal,
		Untagged: untaggedTotal,
	}

	return resp, http.StatusOK, nil
}

func countImageTags(ctx context.Context, b model.Batch, coll *mongo.Collection) (int, int) {
	filter := bson.M{"batch_id": b.ID.Hex()}
	icursor, _ := coll.Find(ctx, filter)
	images := []model.BatchImage{}
	icursor.All(ctx, &images)

	untagged := 0
	for _, img := range images {
		if img.Tag == Untagged {
			untagged += 1
		}
	}
	return len(images) - untagged, untagged
}

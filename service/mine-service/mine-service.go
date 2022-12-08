package mineservice

import (
	"bytes"
	"context"
	"errors"
	"io"
	"path/filepath"
	"time"

	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/constants"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/microservice"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/s3"
	"github.com/workshopapps/pictureminer.api/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DemoMineImage(image io.ReadCloser, filename string) (*model.MineImageResponse, error) {
	content, err := microservice.GetImageContent(image, filename)
	if err != nil {
		return nil, err
	}

	time := time.Now()
	response := &model.MineImageResponse{
		ImageName:    filename,
		TextContent:  content.Content,
		DateCreated:  time,
		DateModified: time,
	}

	return response, nil
}

func MineServiceUpload(userId interface{}, image io.ReadCloser, filename string) (*model.MineImageResponse, error) {
	id, ok := userId.(string)
	if !ok {
		return nil, errors.New("invalid userid")
	}

	image, imageCopy, err := duplicateFile(image)
	if err != nil {
		return nil, err
	}

	imageHash, err := utility.HashFile(imageCopy)
	if err != nil {
		return nil, err
	}

	image, imageCopy, err = duplicateFile(image)
	if err != nil {
		return nil, err
	}

	imagePath, err := s3.UploadImage(image, imageHash+filepath.Ext(filename))
	if err != nil {
		return nil, err
	}

	content, err := microservice.GetImageContent(imageCopy, filename)
	if err != nil {
		return nil, err
	}

	time := time.Now()
	minedImage := &model.MinedImage{
		ID:           primitive.NewObjectID(),
		UserID:       id,
		ImageName:    filename,
		ImageKey:     imageHash,
		ImagePath:    imagePath,
		TextContent:  content.Content,
		DateCreated:  time,
		DateModified: time,
	}

	response, err := saveMinedImage(minedImage, filename)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func GetMinedImages(userId interface{}) ([]model.MineImageResponse, error) {
	id, ok := userId.(string)
	if !ok {
		return nil, errors.New("invalid userid")
	}

	ctx := context.TODO()
	filter := bson.M{"user_id": id}
	cursor, err := mongodb.SelectFromCollection(ctx, config.GetConfig().Mongodb.Database, constants.ImageCollection, filter)
	if err != nil {
		return []model.MineImageResponse{}, err
	}

	minedImages := make([]model.MineImageResponse, 0)
	cursor.All(ctx, &minedImages)

	return minedImages, nil
}

func DeleteMinedImageService(imageKey string) error {
	ctx := context.TODO()
	db := config.GetConfig().Mongodb.Database

	filter := bson.M{"image_key": imageKey}
	minedImageCol := mongodb.GetCollection(mongodb.Connection(), db, constants.ImageCollection)
	_, err := minedImageCol.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func duplicateFile(f io.ReadCloser) (io.ReadCloser, io.ReadCloser, error) {
	contents, err := io.ReadAll(f)
	if err != nil {
		return nil, nil, err
	}
	return io.NopCloser(bytes.NewReader(contents)), io.NopCloser(bytes.NewReader(contents)), nil
}

func saveMinedImage(minedImage *model.MinedImage, filename string) (*model.MineImageResponse, error) {
	_, err := mongodb.MongoPost(constants.ImageCollection, *minedImage)
	if err != nil {
		return nil, err
	}

	// Update API count
	if _, err = mongodb.MongoUpdate(minedImage.UserID[10:len(minedImage.UserID)-2], map[string]interface{}{
		"api_call_count": 1,
	}, constants.UserCollection); err != nil {
		return nil, err
	}

	response := &model.MineImageResponse{
		ImageName:    filename,
		ImageKey:     minedImage.ImageKey,
		ImagePath:    minedImage.ImagePath,
		TextContent:  minedImage.TextContent,
		DateCreated:  minedImage.DateCreated,
		DateModified: minedImage.DateModified,
	}

	return response, nil
}


func ProcessCount(userId string) (model.ProcessCallCount, error) {

	var response model.ProcessCallCount

	batchCount, err := mongodb.MainCountFromCollection(userId, constants.BatchCollection)
	if err != nil {
		return response, err
	}

	imageCount, err := mongodb.MainCountFromCollection(userId, constants.ImageCollection)
	if err != nil {
		return response, err
	}

	var status bool

	if batchCount != 0 || imageCount != 0 {
		status = true
	}else{
		status = false
	}

	response = model.ProcessCallCount{
		ImageCount:    imageCount,
		BatchCount:     batchCount,
		Status:    status,
	}



	return response, nil

}

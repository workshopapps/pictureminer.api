package mineservice

import (
	"bytes"
	"errors"
	"io"
	"path/filepath"
	"time"

	"github.com/workshopapps/pictureminer.api/internal/constants"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/microservice"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/s3"
	"github.com/workshopapps/pictureminer.api/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

	imagePath, err := s3.UploadImage(image, imageHash + filepath.Ext(filename))
	if err != nil {
		return nil, err
	}

	content, err := microservice.GetImageContent(image, filename)
	if err != nil {
		return nil, err
	}

	time := time.Now()
	minedImage := model.MinedImage{
		ID:           primitive.NewObjectID(),
		UserID:       id,
		ImageName:    filename,
		ImageKey:     imageHash,
		ImagePath:    imagePath,
		TextContent:  content.Content,
		DateCreated:  time,
		DateModified: time,
	}

	_, err = mongodb.MongoPost(constants.ImageCollection, minedImage)
	if err != nil {
		return nil, err
	}

	response := &model.MineImageResponse{
		ImageName:    filename,
		ImagePath:    imagePath,
		TextContent:  content.Content,
		DateCreated:  time,
		DateModified: time,
	}

	return response, nil
}

func duplicateFile(f io.ReadCloser) (io.ReadCloser, io.ReadCloser, error) {
	contents, err := io.ReadAll(f)
	if err != nil {
		return nil, nil, err
	}
	return io.NopCloser(bytes.NewReader(contents)), io.NopCloser(bytes.NewReader(contents)), nil
}

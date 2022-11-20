package mineservice

import (
	"io"
	"strings"
	"time"

	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/microservice"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/s3"
	"github.com/workshopapps/pictureminer.api/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	imageCollection = "mined_images"
)

func MineServiceUpload(image io.ReadCloser, filename string) (*model.MineImageResponse, error) {
	hashedImage, err := utility.HashImage(image)
	if err != nil {
		return nil, err
	}

	str := strings.SplitAfter(filename, ".")
	ext := str[len(str)-1]

	imagePath, err := s3.UploadImage(image, hashedImage+"."+ext)
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
		UserID:       "",
		ImageName:    filename,
		ImageKey:     hashedImage,
		ImagePath:    imagePath,
		TextContent:  content.Content,
		DateCreated:  time,
		DateModified: time,
	}

	_, err = mongodb.MongoPost(imageCollection, minedImage)
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

package mineservice

import (
	"mime/multipart"
	"strings"
	"time"

	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/microservice"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongo"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/s3"
	"github.com/workshopapps/pictureminer.api/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	imageCollection = "mined_images"
)

func MineServiceUpload(image multipart.File, filename string) error {
	hashedImage, err := utility.HashImage(image)
	if err != nil {
		return err
	}

	content, err := microservice.GetImageContent(image, filename)
	if err != nil {
		return err
	}

	str := strings.SplitAfter(filename, ".")
	ext := str[len(str)-1]

	imagePath, err := s3.UploadImage(image, hashedImage+"."+ext)
	if err != nil {
		return err
	}

	time := time.Now()
	minedImage := model.MinedImage{
		ID:           primitive.NewObjectID(),
		UserID:       "",
		ImageKey:     hashedImage,
		ImagePath:    imagePath,
		TextContent:  content.Content,
		DateCreated:  time,
		DateModified: time,
	}

	_, err = mongo.MongoPost(imageCollection, minedImage)
	if err != nil {
		return err
	}

	return nil
}

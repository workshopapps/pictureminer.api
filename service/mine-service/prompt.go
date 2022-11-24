package mineservice

import (
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

func PromptMineServiceUpload(userId interface{}, image io.ReadCloser, filename, prompt string) (*model.MineImagePromptResponse, error) {
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

	imagePath, err := s3.UploadImage(image, imageHash+filepath.Ext(filename))
	if err != nil {
		return nil, err
	}

	content, err := microservice.GetImagePromptResponse(image, imageHash, prompt)
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

	response, err := getMineImagePromptResponse(minedImage, content.CheckResult, filename, prompt)
	if err != nil {
		return nil, err
	}

	return response, nil

}

func getMineImagePromptResponse(minedImage *model.MinedImage, checkResult bool, filename, prompt string) (*model.MineImagePromptResponse, error) {
	_, err := mongodb.MongoPost(constants.ImageCollection, *minedImage)
	if err != nil {
		return nil, err
	}

	response := &model.MineImagePromptResponse{
		ImageName:    filename,
		ImagePath:    minedImage.ImagePath,
		TextPrompt:   prompt,
		CheckResult:  checkResult,
		DateCreated:  minedImage.DateCreated,
		DateModified: minedImage.DateModified,
	}

	return response, nil
}

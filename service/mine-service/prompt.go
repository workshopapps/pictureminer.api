package mineservice

import (
	"io"
	"time"

	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/microservice"
	"github.com/workshopapps/pictureminer.api/utility"
)

func PromptMineServiceUpload(image io.ReadCloser, filename, prompt string) (*model.MineImagePromptResponse, error) {
	image, imageCopy, err := duplicateFile(image)
	if err != nil {
		return nil, err
	}

	imageHash, err := utility.HashFile(imageCopy)
	if err != nil {
		return nil, err
	}

	content, err := microservice.GetImagePromptResponse(image, imageHash, prompt)
	if err != nil {
		return nil, err
	}

	time := time.Now()
	response := &model.MineImagePromptResponse{
		ImageName:    filename,
		TextPrompt:   prompt,
		CheckResult:  content.CheckResult,
		DateCreated:  time,
		DateModified: time,
	}

	return response, nil

}

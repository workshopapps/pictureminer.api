package mineservice

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/mocks"
)

func Test_MineServiceUpload(t *testing.T) {
	asst := assert.New(t)
	t.Run("testing upload", func(t *testing.T) {
		testData := model.MineImageResponse{
			ImageName:    "imageName",
			ImagePath:    "ImagePath",
			TextContent:  "TextContent",
			DateCreated:  time.Now(),
			DateModified: time.Now(),
		}
		requestData := model.MineImageUrlRequest{
			Url: "image.jpg",
		}
		mineImageMock := &mocks.MinedImage{TestData: &testData}
		mineImageMock.On("mine image successful", context.Background(), requestData).Return(testData, nil).Once()
		// maniputed the image
		var image io.ReadCloser
		var userId interface{}

		minedImage, err := MineServiceUpload(userId, image, "Filename")
		asst.NoError(err)
		asst.Equal(minedImage, &testData)

		mineImageMock.AssertExpectations(t)

	})
}

func Test_GetMinedImages(t *testing.T) {
	asst := assert.New(t)
	t.Run("testing retrieving mined_image", func(t *testing.T) {
		testData := model.MineImageResponse{
			ImageName:    "imageName",
			ImagePath:    "ImagePath",
			TextContent:  "TextContent",
			DateCreated:  time.Now(),
			DateModified: time.Now(),
		}
		mineImageMock := &mocks.MinedImage{TestData: &testData}
		mineImageMock.On("mine image successful", context.Background()).Return(testData, nil).Once()
		var image io.ReadCloser
		var userId interface{}

		minedImage, err := MineServiceUpload(userId, image, "Filename")
		asst.NoError(err)
		asst.Equal(minedImage, &testData)

		mineImageMock.AssertExpectations(t)
	})
}

package mocks

import (
	"context"
	"io"

	// "github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/workshopapps/pictureminer.api/internal/model"
)

type MinedImage struct {
	mock.Mock
	TestData *model.MineImageResponse
}


func (m *MinedImage) MineServiceUpload(ctx context.Context, url model.MineImageUrlRequest, userId interface{}, image io.ReadCloser, filename string) (*model.MineImageResponse, error) {
	args := m.Called(ctx, url, userId, image, filename)
	return m.TestData, args.Error(1)

}

func (m *MinedImage) GetMinedImages(ctx context.Context, userId interface{}) (*model.MineImageResponse, error) {
	args := m.Called(ctx, userId)
	return m.TestData, args.Error(1)
}

package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	batch "github.com/workshopapps/pictureminer.api/pkg/handler/batch-service"
	"github.com/workshopapps/pictureminer.api/utility"
)

func ProcessBatch(r *gin.Engine, validate *validator.Validate, ApiVersion string, logger *utility.Logger) *gin.Engine {
	batch := batch.Controller{Validate: validate, Logger: logger}

	authUrl := r.Group(fmt.Sprintf("/api/%v", ApiVersion))
	{
		authUrl.POST("/batch-service/process-batch", batch.ProcessBatch)
		authUrl.POST("/batch-service/process-batch-csv", batch.ProcessBatchCSV)
		authUrl.GET("/batch-service/get-batches", batch.GetBatches)
		authUrl.POST("/batch-service/process-batch-api", batch.ProcessBatchAPI)
		authUrl.GET("/batch-service/images/:batch_id", batch.GetBatchImages)
		authUrl.GET("/batch-service/download/:batchid", batch.DownloadCsv)
		authUrl.DELETE("/batch-service/delete/:id", batch.DeleteBatch)
	}

	return r
}

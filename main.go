package main

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/workshopapps/pictureminer.api/internal/config"
	sentry "github.com/workshopapps/pictureminer.api/pkg/middleware"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/s3"
	"github.com/workshopapps/pictureminer.api/utility"

	// "github.com/workshopapps/pictureminer.api/pkg/repository/storage/redis"
	_ "github.com/workshopapps/pictureminer.api/docs"
	"github.com/workshopapps/pictureminer.api/pkg/router"
)

func init() {
	config.Setup()
	// redis.SetupRedis() uncomment when you need redis
	mongodb.ConnectToDB()

	s3.ConnectAws()
	sentry.SentryLogger()
}

// @title           Minergram
// @version         1.0
// @description     A picture mining service API in Go using Gin framework.

// @host      localhost:9000
// @BasePath  /api/v1/
// @schemes http
func main() {
	//Load config
	logger := utility.NewLogger()
	getConfig := config.GetConfig()
	validatorRef := validator.New()
	r := router.Setup(validatorRef, logger)

	logger.Info("Server is starting at 127.0.0.1:%s", getConfig.Server.Port)
	log.Fatal(r.Run(":" + getConfig.Server.Port))
}

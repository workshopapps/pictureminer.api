package main

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/workshopapps/pictureminer.api/internal/config"
	mdb "github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongo"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/s3"
	"github.com/workshopapps/pictureminer.api/utility"

	// "github.com/workshopapps/pictureminer.api/pkg/repository/storage/redis"
	"github.com/workshopapps/pictureminer.api/pkg/router"
)

func init() {
	config.Setup()
	// redis.SetupRedis() uncomment when you need redis
	mdb.ConnectToDB()

	s3.ConnectAws()
}

func main() {
	//Load config
	logger := utility.NewLogger()
	getConfig := config.GetConfig()
	validatorRef := validator.New()
	r := router.Setup(validatorRef, logger)

	logger.Info("Server is starting at 127.0.0.1:%s", getConfig.Server.Port)
	log.Fatal(r.Run(":" + getConfig.Server.Port))
}

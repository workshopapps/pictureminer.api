package main

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/workshopapps/pictureminer.api/internal/config"
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
}

// @title           Minergram
// @version         1.0
// @description     A picture mining service API in Go using Gin framework.

// @host      discripto.hng.tech/api1
// @BasePath  /api/v1/
// @schemes https
// @securityDefinitions.apikey BearerAuth
// @in header "Bearer <add access token here>"
// @name Authorization
func main() {
	// err := sentry.Init(sentry.ClientOptions{
	// 	Dsn: "https://419447b5b02e42dc8b277f5af67e565f@o4504279417421824.ingest.sentry.io/4504279420305408",
	// 	// Set TracesSampleRate to 1.0 to capture 100%
	// 	// of transactions for performance monitoring.
	// 	// We recommend adjusting this value in production,
	// 	TracesSampleRate: 1.0,
	// 	// TracesSampler: sentry.TracesSamplerFunc(func(ctx sentry.SamplingContext) sentry.Sampled {
	// 	// 	return sentry.SampledTrue
	// 	// }),
	// })
	// if err != nil {
	// 	log.Fatalf("sentry.Init: %s", err)
	// }
	// // Flush buffered events before the program terminates.
	// defer sentry.Flush(2 * time.Second)
	// sentry.CaptureException(err)
	// sentry.CaptureMessage("It works fine!")
	//Load config
	logger := utility.NewLogger()
	getConfig := config.GetConfig()
	validatorRef := validator.New()
	r := router.Setup(validatorRef, logger)

	logger.Info("Server is starting at 127.0.0.1:%s", getConfig.Server.Port)
	log.Fatal(r.Run(":" + getConfig.Server.Port))

}

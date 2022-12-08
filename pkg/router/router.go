package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/workshopapps/pictureminer.api/pkg/middleware"
	"github.com/workshopapps/pictureminer.api/utility"
)

func Setup(validate *validator.Validate, logger *utility.Logger) *gin.Engine {

	// To initialize Sentry's handler, you need to initialize Sentry itself beforehand
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://419447b5b02e42dc8b277f5af67e565f@o4504279417421824.ingest.sentry.io/4504279420305408",
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

	sentry.CaptureMessage("It works fine o!")

	r := gin.New()

	r.Use(sentrygin.New(sentrygin.Options{}))

	r.Use(func(ctx *gin.Context) {
		if hub := sentrygin.GetHubFromContext(ctx); hub != nil {
			hub.Scope().SetTag("someRandomTag", "maybeYouNeedIt")
		}
		ctx.Next()
	})

	r.GET("/testing", func(ctx *gin.Context) {
		if hub := sentrygin.GetHubFromContext(ctx); hub != nil {
			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetExtra("unwantedQuery", "someQueryDataMaybe")
				hub.CaptureMessage("User provided unwanted query string, but we recovered just fine")
			})
		}
		ctx.Status(http.StatusOK)
	})

	// Middlewares
	// r.Use(gin.Logger())

	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.MaxMultipartMemory = 1 << 20 // 1MB

	ApiVersion := "v1"
	Health(r, validate, ApiVersion, logger)
	Auth(r, validate, ApiVersion, logger)
	ProcessBatch(r, validate, ApiVersion, logger)
	MineService(r, validate, ApiVersion, logger)
	Admin(r, validate, ApiVersion, logger)
	SwaggerDocs(r, ApiVersion)
	Feedback(r, validate, ApiVersion, logger)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"name":    "Not Found",
			"message": "Page not found.",
			"code":    404,
			"status":  http.StatusNotFound,
		})
	})

	return r
}

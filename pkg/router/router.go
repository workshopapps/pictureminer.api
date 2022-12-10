package router

import (
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.atatus.com/agent/module/atgin"

	// newrelic "github.com/newrelic/go-agent"
	// "github.com/newrelic/go-agent/_integrations/nrgin/v1"
	"github.com/workshopapps/pictureminer.api/pkg/middleware"
	"github.com/workshopapps/pictureminer.api/utility"
)

func Setup(validate *validator.Validate, logger *utility.Logger) *gin.Engine {

	// cfg := newrelic.NewConfig("discripto_api", "23e1bbb04e4fd6b88bdedb97fde89345ee8cNRAL")

	// app, err := newrelic.NewApplication(cfg)
	// if nil != err {
	// 	fmt.Println(err)
	// }

	r := gin.New()
	r.Use(atgin.Middleware(r))
	// r.Use(nrgin.Middleware(app))

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

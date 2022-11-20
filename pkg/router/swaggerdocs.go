package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SwaggerDocs(r *gin.Engine,  ApiVersion string) *gin.Engine {
	docsUrl := r.Group(fmt.Sprintf("/api/%v", ApiVersion))
	{
		docsUrl.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	return r
}

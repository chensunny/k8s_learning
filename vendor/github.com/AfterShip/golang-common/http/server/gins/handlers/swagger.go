package handlers

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/AfterShip/golang-common/logger"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/swaggo/swag"
	"go.uber.org/zap"
)

type SwaggerDoc struct {
	swaggerDocFilePath string
}

func registerSwaggerDoc(swaggerDocFilePath string) {
	swag.Register(swag.Name, &SwaggerDoc{
		swaggerDocFilePath: swaggerDocFilePath,
	})
}

func (doc *SwaggerDoc) ReadDoc() string {
	content, err := ioutil.ReadFile(doc.swaggerDocFilePath)
	if err != nil {
		logger.Info(context.TODO(), "We cannot load swagger file.", zap.Error(err))
	}
	return string(content)
}

func RegisterSwaggerDocHandler(group *gin.RouterGroup, swaggerFilePath string) {
	registerSwaggerDoc(swaggerFilePath)
	group.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	group.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "swagger/index.html")
	})
}

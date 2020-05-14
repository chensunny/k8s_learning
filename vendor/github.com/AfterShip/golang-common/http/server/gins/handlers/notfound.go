package handlers

import (
	"github.com/AfterShip/golang-common/errors"
	"github.com/AfterShip/golang-common/http/server/gins"
	"github.com/gin-gonic/gin"
)

func NoRoute(c *gin.Context) {
	gins.ResponseError(c, nil, errors.ErrNotFound)
}

func NoMethod(c *gin.Context) {
	gins.ResponseError(c, nil, errors.ErrMethodNotAllowed)
}

func RegisterNotFoundHandlers(engine *gin.Engine) {
	engine.NoRoute(NoRoute)
	engine.NoMethod(NoMethod)
}

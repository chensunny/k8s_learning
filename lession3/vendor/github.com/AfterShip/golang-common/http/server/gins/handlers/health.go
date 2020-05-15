package handlers

import (
	"github.com/AfterShip/golang-common/http/server/health"
	"github.com/gin-gonic/gin"
)

func RegisterHealthHandler(engine *gin.Engine, s *health.Status) {
	statusHandler := gin.WrapF(s.HttpHandlerFunc())
	engine.PUT(s.RequestURI(), statusHandler)
	engine.POST(s.RequestURI(), statusHandler)
	engine.GET(s.RequestURI(), statusHandler)
	engine.HEAD(s.RequestURI(), statusHandler)
}

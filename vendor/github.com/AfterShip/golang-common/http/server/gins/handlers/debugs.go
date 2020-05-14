package handlers

import (
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/trace"
)

func RegisterDebugHandler(engine *gin.Engine) {
	// debug group
	debugGroup := engine.Group("debug")

	debugGroup.GET("requests", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		trace.Render(c.Writer, c.Request, true)
	})
	debugGroup.GET("events", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		trace.RenderEvents(c.Writer, c.Request, true)
	})

	// pprof sub group
	pprofGroup := debugGroup.Group("pprof")
	pprofGroup.GET("/", gin.WrapF(pprof.Index))
	pprofGroup.GET("/:name", func(c *gin.Context) {
		var handlerFunc http.HandlerFunc
		name := c.Param("name")
		switch name {
		case "cmdline":
			handlerFunc = pprof.Cmdline
		case "profile":
			handlerFunc = pprof.Profile
		case "symbol":
			handlerFunc = pprof.Symbol
		case "trace":
			handlerFunc = pprof.Trace
		default:
			handlerFunc = pprof.Handler(name).ServeHTTP
		}

		handlerFunc(c.Writer, c.Request)
	})
}

package main

import (
	"fmt"
	"net/http"

	"github.com/AfterShip/golang-common/http/server/gins/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	devopsHttpServer := &http.Server{
		Addr: fmt.Sprintf(":%d", 8080),
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	//debug info
	//paths: /debug/requests、/debug/events、/debug/pprof
	handlers.RegisterDebugHandler(engine)

	//whoami
	handlers.RegisterWhoamiHandler(engine.Group(""))

	devopsHttpServer.Handler = engine

	devopsHttpServer.ListenAndServe()
}

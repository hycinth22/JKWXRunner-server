package apiServer

import (
	"github.com/gin-gonic/gin"
	"github.com/inkedawn/JKWXFucker-server/config"
)

func Run(engine *gin.Engine) error {
	registerRoute(engine, RouteRegisterFunctions...)
	engine.Use(cacheControlMiddleWare)
	return engine.Run(config.ListenAddr)
}

type RouteRegisterFunc func(router gin.IRouter)

func registerRoute(engine *gin.Engine, allF ...RouteRegisterFunc) {
	for _, f := range allF {
		f(engine)
	}
}

func cacheControlMiddleWare(c *gin.Context) {
	c.Header("Cache-Control", "no-store, max-age=0")
}

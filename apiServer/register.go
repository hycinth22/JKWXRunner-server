package apiServer

import "github.com/gin-gonic/gin"

type RequestRouter func(router gin.IRouter)

// register action only happened at call Server.Run()
var RequestRoutersTable []RequestRouter = nil

func registerRequestRoutersTable(engine *gin.Engine) {
	registerRouterRule(engine, RequestRoutersTable...)
}

func registerRouterRule(engine *gin.Engine, allF ...RequestRouter) {
	for _, f := range allF {
		f(engine)
	}
}

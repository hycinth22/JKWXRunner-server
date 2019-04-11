package apiServer

import "github.com/gin-gonic/gin"

type RequestRouter func(router gin.IRouter)

// read-only, register action only happened at call Run()
var RequestRoutersTable = []RequestRouter{
	CORSRouter,
	AccountRouter,
	AccountLogsRouter,
	UserInfoRouter,
}

func registerRequestRoutersTable(engine *gin.Engine) {
	registerRouterRule(engine, RequestRoutersTable...)
}

func registerRouterRule(engine *gin.Engine, allF ...RequestRouter) {
	for _, f := range allF {
		f(engine)
	}
}

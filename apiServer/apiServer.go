package apiServer

import "github.com/gin-gonic/gin"

func Run(engine *gin.Engine) error {
	registerCORSRoute(engine)
	registerTicketRoute(engine)
	registerRemoteProfileRoute(engine)
	engine.Use(cacheMiddleWare)
	return engine.Run(":8080")
}

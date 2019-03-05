package apiServer

import (
	"github.com/gin-gonic/gin"
	"github.com/inkedawn/JKWXFucker-server/config"
)

func Run(engine *gin.Engine) error {
	registerCORSRoute(engine)
	registerTicketRoute(engine)
	registerRemoteProfileRoute(engine)
	engine.Use(cacheMiddleWare)
	return engine.Run(config.ListenAddr)
}

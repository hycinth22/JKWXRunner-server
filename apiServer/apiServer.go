package apiServer

import "github.com/gin-gonic/gin"

func Run(engine *gin.Engine) error {
	registerCORSRoute(engine)
	registerTicketRoute(engine)
	registerRemoteProfileRoute(engine)
	return engine.Run(":8080")
}

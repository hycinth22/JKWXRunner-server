package apiServer

import (
	"github.com/gin-gonic/gin"
	"github.com/inkedawn/JKWXRunner-server/config"
)

func Run(engine *gin.Engine) error {
	registerRequestRoutersTable(engine)
	engine.Use(addCacheControlHeader)
	return engine.Run(config.ListenAddr)
}

func addCacheControlHeader(c *gin.Context) {
	c.Header("Cache-Control", "no-store, max-age=0")
}

// DEPRECATED
package staticFileServer

import "github.com/gin-gonic/gin"

// DEPRECATED
func Run(engine *gin.Engine) error {
	// static files
	engine.Use(authMiddleWare)
	engine.Use(cacheMiddleWare)
	engine.Static("/", `./html`)
	return engine.Run(":80")
}

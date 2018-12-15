package staticFileServer

import (
	"github.com/gin-gonic/gin"
)

func cacheMiddleWare(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, public")
}

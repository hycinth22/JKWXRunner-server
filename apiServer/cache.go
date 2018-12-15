package apiServer

import (
	"github.com/gin-gonic/gin"
)

func cacheMiddleWare(c *gin.Context) {
	c.Header("Cache-Control", "no-store, max-age=0")
}

package apiServer

import (
	"github.com/gin-gonic/gin"
)

func registerCORSRoute(router gin.IRouter) {
	router.Use(func(context *gin.Context) {
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")

	})
	router.OPTIONS("/:all", func(context *gin.Context) {
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
		context.Header("Access-Control-Allow-Headers", "*")
	})
}

package apiServer

import (
	"github.com/gin-gonic/gin"
	"github.com/inkedawn/JKWXFucker-server/apiServer/handler"
)

func CORSRouter(router gin.IRouter) {
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

func AccountRouter(router gin.IRouter) {
	router.GET("/account", handler.ListAccounts)
	router.POST("/account", handler.NotImplemented)
	router.PUT("/account", handler.NotImplemented)
	router.DELETE("/account/:id", handler.NotImplemented)
}

func AccountLogsRouter(router gin.IRouter) {
	router.GET("/account/:uid/logs/", handler.ListLogsByUID)
}

func UserInfoRouter(router gin.IRouter) {
	router.GET("/userInfo/:username", handler.NotImplemented)
}

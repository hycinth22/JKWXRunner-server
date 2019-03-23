package apiServer

import (
	"github.com/gin-gonic/gin"
	"github.com/inkedawn/JKWXFucker-server/apiServer/handler"
)

// read-only, register action only happened at call Run()
var RouteRegisterFunctions = []RouteRegisterFunc{
	registerCORSRoute,
	registerAccountRoute,
	registerLogRoute,
}

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

func registerAccountRoute(router gin.IRouter) {
	router.GET("/account", handler.ListAccounts)
	router.POST("/account", handler.AddAccount)
	router.PUT("/account", handler.UpdateAccount)
	// router.DELETE("/account/:id", handler.DeleteAccount)
}

func registerLogRoute(router gin.IRouter) {
	router.GET("/log", handler.ListLogsByUsername)
}

func registerUserInfoRoute(router gin.IRouter) {
	router.GET("/userInfo/:username", handler.ShowUserInfo)
}
